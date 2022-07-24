package app

import (
	"github.com/mikelorant/ezdb2/internal/storage"
	"github.com/mikelorant/ezdb2/internal/structprinter"
)

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

func (a *App) GetStorageClient(name string) (*storage.Store, error) {
	store := a.Config.getStore(name)

	cfg := storage.Config{
		Name:   store.Name,
		Type:   store.Type,
		Path:   store.Path,
		Bucket: store.Bucket,
		Prefix: store.Prefix,
		Region: store.Region,
	}

	return storage.New(cfg)
}
