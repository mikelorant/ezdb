package compress

import (
	"compress/gzip"
	"io"
)

// source (r) -> (r) copy (w) -> (w) gzip (w) -> (w) pipe (r)
func NewGzipCompressor(source io.Reader) io.Reader {
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

func NewGzipDecompressor(rc io.Reader) io.Reader {
	r, err := gzip.NewReader(rc)
	if err != nil {
		r.Close()
	}
	return r
}
