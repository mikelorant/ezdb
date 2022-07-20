package storage

import (
	"fmt"
	"io"
)

type PipeStore struct {
	reader *io.PipeReader
	writer *io.PipeWriter
}

func NewPipeStorer() (*PipeStore, error) {
	r, w := io.Pipe()

	return &PipeStore{
		reader: r,
		writer: w,
	}, nil
}

func (p *PipeStore) Store(data io.Reader, filename string, done chan bool, result chan string) error {
	go func() error {
		defer p.writer.Close()

		_, err := io.Copy(p.writer, data)
		if err != nil {
			return fmt.Errorf("unable to write to pipe: %w", err)
		}

		result <- ""
		done <- true

		return nil
	}()

	return nil
}

func (p *PipeStore) Retrieve(data io.WriteCloser, filename string, done chan bool) error {
	go func() error {
		_, err := io.Copy(data, p.reader)
		if err != nil {
			return fmt.Errorf("unable to read from pipe: %w", err)
		}
		data.Close()
		done <- true

		return nil
	}()

	return nil
}

func (p *PipeStore) List() ([]string, error) {
	return []string{}, nil
}
