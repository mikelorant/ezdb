package storage

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"
)

type FileStore struct {
	Directory string
}

type FileOptions struct {
	Data     io.Reader
	Filename string
}

const (
	FilenameFormat = "%v-20060102-150405.sql.gz"
)

func NewFileStorer(directory string) (*FileStore, error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, fmt.Errorf("Error mkdir: %v: %w", directory, err)
	}

	return &FileStore{
		Directory: directory,
	}, nil
}

func (f *FileStore) Store(data io.Reader, filename string, done chan bool, result chan string) error {
	filename = fmt.Sprintf(FilenameFormat, filename)
	filename = time.Now().Format(filename)
	filename = path.Join(f.Directory, filename)

	fd, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create file: %w", err)
	}

	go func() error {
		defer fd.Close()

		_, err := io.Copy(fd, data)
		if err != nil {
			return fmt.Errorf("unable to write file: %w", err)
		}

		result <- fd.Name()
		done <- true

		return nil
	}()

	return nil
}

func (f *FileStore) Retrieve(data io.WriteCloser, filename string, done chan bool) error {
	filename = path.Join(f.Directory, filename)

	fd, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}

	go func() error {
		defer fd.Close()

		_, err := io.Copy(data, fd)
		if err != nil {
			return fmt.Errorf("unable to read file: %w", err)
		}
		data.Close()

		done <- true

		return nil
	}()

	return nil
}
