package app

import (
	"io"

	"github.com/mikelorant/ezdb2/internal/structprinter"
)

type Storer interface {
	Store(data io.Reader, filename string, done chan bool, result chan string) error
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
