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

func (p *PipeStore) Store(data io.Reader, filename string) (string, error) {
	defer p.writer.Close()

	_, err := io.Copy(p.writer, data)
	if err != nil {
		return "", fmt.Errorf("unable to write to pipe: %w", err)
	}

	return "", nil
}

func (p *PipeStore) Retrieve(data io.WriteCloser, filename string) error {
	_, err := io.Copy(data, p.reader)
	if err != nil {
		return fmt.Errorf("unable to read from pipe: %w", err)
	}
	data.Close()

	return nil
}

func (p *PipeStore) List() ([]string, error) {
	return []string{}, nil
}
