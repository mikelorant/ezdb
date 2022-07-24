package mysqlshim

import (
	"fmt"
)

const (
	QueryCreateDatabase = "CREATE DATABASE IF NOT EXISTS %v;"
)

func (s *Shim) CreateDatabase(name string) error {
	q := fmt.Sprintf(QueryCreateDatabase, name)

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	_, err = tx.Exec(q)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}
	tx.Commit()

	return nil
}
