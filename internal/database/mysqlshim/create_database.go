package mysqlshim

import (
	"fmt"
)

const (
	QueryCreateDatabase = "CREATE DATABASE IF NOT EXISTS %v;"
)

func (s *Shim) CreateDatabase(name string) error {
	q := fmt.Sprintf(QueryCreateDatabase, name)

	if err := s.exec(q); err != nil {
		return fmt.Errorf("unable to create database: %v: %w", name, err)
	}

	return nil
}
