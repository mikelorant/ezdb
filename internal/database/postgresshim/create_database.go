package postgresshim

import (
	"fmt"
)

const (
	QueryCreateDatabase = "CREATE DATABASE %v;"
)

func (s *Shim) CreateDatabase(name string) error {
	q := fmt.Sprintf(QueryCreateDatabase, name)

	if _, err := s.DB.Exec(q); err != nil {
		return fmt.Errorf("unable to exec query: %w", err)
	}

	return nil
}
