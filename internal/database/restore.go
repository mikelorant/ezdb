package database

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

type ReplaceRegexpString [2]string

var (
	MySQLRestoreReplaceUTF     = ReplaceRegexpString{"utf8mb4_0900_ai_ci", "utf8mb4_unicode_ci"}
	MySQLRestoreReplaceDefiner = ReplaceRegexpString{"DEFINER=[^ *]+", "DEFINER=CURRENT_USER"}
)

func (d *Database) Restore(name, filename string, storer Storer, shell Shell, verbose bool) ([]byte, error) {
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
		return storer.Retrieve(pw, filename)
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
	))

	if verbose {
		log.Println("Command:", d.RestoreCommand(true))
	}

	if err := shell.Run(&buf, rr, d.RestoreCommand(false), true); err != nil {
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

func getTransformer(r ReplaceRegexpString) transform.Transformer {
	return replace.RegexpString(regexp.MustCompile(r[0]), r[1])
}
