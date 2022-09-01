package app

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"

	"github.com/icholy/replace"
	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
	"golang.org/x/sync/errgroup"
	"golang.org/x/text/transform"
)

type Retriever interface {
	Retrieve(data io.WriteCloser, filename string) error
}

type Restorer interface {
	RestoreCommand(verbose bool) string
}

type RestoreOptions struct {
	Context  string
	Name     string
	Store    string
	Filename string
}

func (a *App) Restore(opts RestoreOptions) error {
	context, err := Select(opts.Context, a.Config.getContexts(), "Choose a context:")
	if err != nil {
		return fmt.Errorf("unable to select a context: %w", err)
	}

	cl, err := a.GetDBClient(context)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	store, err := Select(opts.Store, a.Config.getStores(), "Choose a store:")
	if err != nil {
		return fmt.Errorf("unable to select a store: %w", err)
	}

	storer, err := a.GetStorageClient(store)
	if err != nil {
		return fmt.Errorf("unable to get storer: %w", err)
	}

	filenames, err := storer.List()
	if err != nil {
		return fmt.Errorf("unable to list store: %w", err)
	}

	filename, err := Select(opts.Filename, filenames, "Choose a file:")
	if err != nil {
		return fmt.Errorf("unable to select a file: %w", err)
	}

	if cl.IsDatabase(opts.Name) {
		log.Println("Found existing database:", opts.Name)
		log.Println("Press enter to drop database.")
		fmt.Scanln()

		if err := cl.DropDatabase(opts.Name); err != nil {
			return fmt.Errorf("unable to drop target database: %v: %w", opts.Name, err)
		}
	}

	if err := cl.CreateDatabase(opts.Name); err != nil {
		return fmt.Errorf("unable to create database: %v: %w", opts.Name, err)
	}

	cl, err = a.GetDBClient(context,
		WithDBName(opts.Name),
	)
	if err != nil {
		return fmt.Errorf("unable to get database client: %w", err)
	}

	shell, err := a.GetShell(context)
	if err != nil {
		return fmt.Errorf("unable to get a shell: %w", err)
	}

	_, err = doRestore(cl, opts.Name, filename, storer, shell, true)
	if err != nil {
		return fmt.Errorf("unable to restore database: %v: %w", opts.Name, err)
	}

	log.Println("Database successfully restored.")

	return nil
}

type ReplaceRegexpString [2]string

var (
	MySQLRestoreReplaceUTF     = ReplaceRegexpString{"utf8mb4_0900_ai_ci", "utf8mb4_unicode_520_ci"}
	MySQLRestoreReplaceDefiner = ReplaceRegexpString{"DEFINER=[^ *]+", "DEFINER=CURRENT_USER"}
	MySQLRestoreReplaceGTID    = ReplaceRegexpString{"SET\\ @@GLOBAL.GTID_PURGED=", "#\\ SET\\ @@GLOBAL.GTID_PURGED="}
)

func doRestore(cmd Restorer, name, filename string, retriever Retriever, runner Runner, verbose bool) ([]byte, error) {
	desc := "Restoring..."
	bar := progress.New(-1, desc, verbose)

	var buf bytes.Buffer

	// storer (w) -> (w) pipe (r) -> (r) teereader (w) -> (w) progressbar
	// 			     			                   (r) -> (r) buffer (r) -> (r) gzip (r) -> (r) replacer (r) -> (r) mysql

	pr, pw := io.Pipe()
	tr := io.TeeReader(pr, bar)
	rb := bufio.NewReader(tr)

	g := new(errgroup.Group)

	g.Go(func() error {
		return retriever.Retrieve(pw, filename)
	})

	var r io.Reader
	r = rb

	gz, err := rb.Peek(2)
	if err != nil {
		return nil, fmt.Errorf("unable to check for compression: %w", err)
	}
	if gz[0] == 31 && gz[1] == 139 {
		rs := compress.NewGzipDecompressor(rb)
		r = rs
	}

	rr := transform.NewReader(r, transform.Chain(
		getTransformer(MySQLRestoreReplaceUTF),
		getTransformer(MySQLRestoreReplaceDefiner),
		getTransformer(MySQLRestoreReplaceGTID),
	))

	if verbose {
		log.Println("Command:", cmd.RestoreCommand(true))
	}

	if err := runner.Run(&buf, rr, cmd.RestoreCommand(false), true); err != nil {
		bar.Finish()
		out, _ := io.ReadAll(&buf)
		fmt.Print(string(out))
		return nil, fmt.Errorf("unable to run command: %w", err)
	}

	out, err := io.ReadAll(&buf)
	if err != nil {
		return nil, fmt.Errorf("unable to read output: %w", err)
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("store failure: %w", err)
	}

	bar.Finish()

	return out, nil
}

func getTransformer(r ReplaceRegexpString) *replace.RegexpTransformer {
	return replace.RegexpString(regexp.MustCompile(r[0]), r[1])
}
