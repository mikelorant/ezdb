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

func NewFileStore(directory string) (*FileStore, error) {
	if err := os.MkdirAll(directory, 0o755); err != nil {
		return nil, fmt.Errorf("Error mkdir: %v: %w", directory, err)
	}

	return &FileStore{
		Directory: directory,
	}, nil
}

func (f *FileStore) Store(data io.Reader, filename string) (string, error) {
	filename = fmt.Sprintf(FilenameFormat, filename)
	filename = time.Now().Format(filename)
	filename = path.Join(f.Directory, filename)

	fd, err := os.Create(filename)
	if err != nil {
		return "", fmt.Errorf("unable to create file: %w", err)
	}
	defer fd.Close()

	_, err = io.Copy(fd, data)
	if err != nil {
		return "", fmt.Errorf("unable to write file: %w", err)
	}

	return fd.Name(), nil
}

func (f *FileStore) Retrieve(data io.WriteCloser, filename string) error {
	filename = path.Join(f.Directory, filename)

	fd, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file: %w", err)
	}
	defer fd.Close()

	_, err = io.Copy(data, fd)
	if err != nil {
		return fmt.Errorf("unable to read file: %w", err)
	}
	data.Close()

	return nil
}

func (f *FileStore) List() ([]string, error) {
	var list []string

	files, err := os.ReadDir(f.Directory)
	if err != nil {
		return list, fmt.Errorf("unable to read directory: %v: %w", f.Directory, err)
	}

	for _, file := range files {
		list = append(list, file.Name())
	}

	return list, nil
}
