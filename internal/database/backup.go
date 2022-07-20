package database

import (
	"compress/gzip"
	"database/sql"
	"fmt"
	"io"

	"github.com/mikelorant/ezdb2/internal/progress"

	"github.com/jamf/go-mysqldump"
)

func (cl *Client) Backup(name string, size int64, storer Storer) (string, error) {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Dumping..."
	bar := progress.New(size, desc)

	done := make(chan bool)
	result := make(chan string)

	// mysqldump (w) -> (w) multiwriter (w) -> (w) progressbar
	//                                      -> (w) pipe (r) -> (r) gzip (r) -> (r) storer

	gzipIn, dumpOut := io.Pipe()
	dumpIn := io.MultiWriter(dumpOut, bar)
	gzipOut := NewGzipReader(gzipIn)

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

// source (r) -> (r) copy (w) -> (w) gzip (w) -> (w) pipe (r)
func NewGzipReader(source io.Reader) io.Reader {
	r, w := io.Pipe()
	go func() {
		defer w.Close()

		zip, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		defer zip.Close()
		if err != nil {
			w.CloseWithError(err)
		}

		io.Copy(zip, source)
	}()
	return r
}
