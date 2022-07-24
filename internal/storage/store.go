package storage

import (
	"fmt"
	"io"

	"github.com/mikelorant/ezdb2/internal/storage/bucketstore"
	"github.com/mikelorant/ezdb2/internal/storage/filestore"
	"github.com/mikelorant/ezdb2/internal/storage/pipestore"
)

type Storer interface {
	Store(data io.Reader, filename string) (string, error)
	Retrieve(data io.WriteCloser, filename string) error
	List() ([]string, error)
}

type Config struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Path   string `yaml:"path"`
	Bucket string `yaml:"bucket"`
	Prefix string `yaml:"prefix"`
	Region string `yaml:"region"`
}

type Store struct {
	Config Config
	storer Storer
}

func New(cfg Config) (*Store, error) {
	switch cfg.Type {
	case "s3":
		s, err := bucketstore.New(cfg.Region, cfg.Bucket, cfg.Prefix)
		if err != nil {
			return nil, fmt.Errorf("unable to provision storage: %v: %w", cfg.Name, err)
		}
		return &Store{
			Config: cfg,
			storer: s,
		}, nil
	case "directory":
		s, err := filestore.New(cfg.Path)
		if err != nil {
			return nil, fmt.Errorf("unable to provision storage: %v: %w", cfg.Name, err)
		}
		return &Store{
			Config: cfg,
			storer: s,
		}, nil
	case "pipe":
		s, err := pipestore.New()
		if err != nil {
			return nil, fmt.Errorf("unable to provision storage: %v: %w", cfg.Name, err)
		}
		return &Store{
			Config: cfg,
			storer: s,
		}, nil
	}

	return &Store{}, nil
}

func (s *Store) Store(data io.Reader, filename string) (string, error) {
	return s.storer.Store(data, filename)
}

func (s *Store) Retrieve(data io.WriteCloser, filename string) error {
	return s.storer.Retrieve(data, filename)
}

func (s *Store) List() ([]string, error) {
	return s.storer.List()
}
