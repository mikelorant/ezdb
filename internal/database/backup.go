package database

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/jamf/go-mysqldump"
	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
	"golang.org/x/sync/errgroup"
)

type Shell interface {
	Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error
}

var (
	MySQLDumpCommand = "mysqldump"
	MySQLDumpOptions = []string{
		"--compress",
		"--column-statistics=0",
		"--ssl-mode=preferred",
		"--set-gtid-purged=OFF",
	}
)

func (cl *Client) Backup(name string, size int64, storer Storer, verbose bool) (string, error) {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Dumping..."
	bar := progress.New(size, desc, verbose)

	done := make(chan bool)
	result := make(chan string)

	// mysqldump (w) -> (w) multiwriter (w) -> (w) progressbar
	//                                      -> (w) pipe (r) -> (r) gzip (r) -> (r) storer

	gzipIn, dumpOut := io.Pipe()
	dumpIn := io.MultiWriter(dumpOut, bar)
	gzipOut := compress.NewGzipCompressor(gzipIn)

	storer.Store(gzipOut, name)

	dumper := &mysqldump.Data{
		Connection: db,
		Out:        dumpIn,
	}
	if err := dumper.Dump(); err != nil {
		return "", fmt.Errorf("unable to dump database: %w", err)
	}
	dumpOut.Close()
	location := <-result
	<-done
	bar.Finish()

	return location, nil
}

func (cl *Client) BackupCompat(name string, size int64, storer Storer, shell Shell, verbose bool) (string, error) {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Dumping..."
	bar := progress.New(size, desc, verbose)

	// mysqldump (w) -> (w) multiwriter (w) -> (w) progressbar
	//                                      -> (w) pipe (r) -> (r) gzip (r) -> (r) storer

	gzipIn, dumpOut := io.Pipe()
	dumpIn := io.MultiWriter(dumpOut, bar)
	gzipOut := compress.NewGzipCompressor(gzipIn)

	g := new(errgroup.Group)

	result := make(chan string)

	g.Go(func() error {
		location, err := storer.Store(gzipOut, name)
		result <- location
		return err
	})

	log.Println("Command:", cl.getBackupCommand(true))
	if err := shell.Run(dumpIn, nil, cl.getBackupCommand(false), false); err != nil {
		return "", fmt.Errorf("unable to run command: %w", err)
	}
	dumpOut.Close()

	location := <-result
	if err := g.Wait(); err != nil {
		return "", fmt.Errorf("store failure: %w", err)
	}
	bar.Finish()

	return location, nil
}

func (cl *Client) getBackupCommand(hidden bool) string {
	var cmd []string

	cmd = append(cmd, MySQLDumpCommand)

	if user := cl.config.User; user != "" {
		cmd = append(cmd, fmt.Sprintf("--user=%v", user))
	}

	if password := cl.config.Passwd; password != "" {
		if hidden {
			cmd = append(cmd, fmt.Sprintf("--password=%v", strings.Repeat("*", len(password))))
		} else {
			cmd = append(cmd, fmt.Sprintf("--password=%v", password))
		}
	}

	hostPort := strings.Split(cl.config.Addr, ":")
	cmd = append(cmd, fmt.Sprintf("--host=%v --port=%v", hostPort[0], hostPort[1]))

	cmd = append(cmd, MySQLDumpOptions...)

	cmd = append(cmd, cl.config.DBName)

	return strings.Join(cmd, " ")
}
