package database

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/jamf/go-mysqldump"
	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
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

	storer.Store(gzipOut, name, done, result)

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
