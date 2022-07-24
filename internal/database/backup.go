package database

import (
	"fmt"
	"io"
	"log"

	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
	"golang.org/x/sync/errgroup"
)

func (d *Database) Backup(name string, size int64, storer Storer, shell Shell, verbose bool) (string, error) {
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

	if verbose {
		log.Println("Command:", d.BackupCommand(true))
	}

	if err := shell.Run(dumpIn, nil, d.BackupCommand(false), false); err != nil {
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
