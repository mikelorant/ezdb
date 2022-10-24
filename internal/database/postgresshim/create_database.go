package postgresshim

import (
	"fmt"
)

const (
	QueryCreateDatabase = "CREATE DATABASE %v;"
)

func (s *Shim) CreateDatabase(name string) error {
	q := fmt.Sprintf(QueryCreateDatabase, name)

	if _, err := s.query(q); err != nil {
		return fmt.Errorf("unable to create database: %v: %w", name, err)
	}

	return nil
}
