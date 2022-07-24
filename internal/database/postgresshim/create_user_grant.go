package postgresshim

import (
	"fmt"
)

func (s *Shim) CreateUserGrant(name, password, database string) error {
	query := []string{
		fmt.Sprintf("CREATE USER %v WITH PASSWORD '%v';", name, password),
		fmt.Sprintf("GRANT ALL ON DATABASE %v TO %v;", database, name),
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	for _, q := range query {
		_, err = tx.Exec(q)
		if err != nil {
			return fmt.Errorf("unable to begin transaction: %w", err)
		}
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("transaction failed: %w", err)
		}
	}
	tx.Commit()

	return nil
}
