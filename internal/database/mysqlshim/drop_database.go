package mysqlshim

import (
	"fmt"
)

func (s *Shim) DropDatabase(name string) error {
	q := fmt.Sprintf("DROP DATABASE IF EXISTS %v;", name)

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
