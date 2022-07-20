package app

import (
	"fmt"
	"io"

	"github.com/mikelorant/ezdb2/internal/storage"
	"github.com/mikelorant/ezdb2/internal/structprinter"
)

type Storer interface {
	Store(data io.Reader, filename string, done chan bool, result chan string) error
	Retrieve(data io.WriteCloser, filename string, done chan bool) error
	List() ([]string, error)
}

type Stores []Store

type Store struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Bucket string `yaml:"bucket"`
	Prefix string `yaml:"prefix"`
	Region string `yaml:"region"`
}

func (s Store) String() string {
	return structprinter.Sprint(s)
}

func GetStorer(store *Store) (Storer, error) {
	var storer Storer
	switch store.Type {
	case "s3":
		s, err := storage.NewBucketStorer(store.Region, store.Bucket, store.Prefix)
		if err != nil {
			return nil, fmt.Errorf("unable to provision storage: %v: %w", store.Name, err)
		}
		storer = s
	case "directory":
		s, err := storage.NewFileStorer(store.Path)
		if err != nil {
			return nil, fmt.Errorf("unable to provision storage: %v: %w", store.Name, err)
		}
		storer = s
	}

	return storer, nil
}
