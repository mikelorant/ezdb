package mysqlshim

import (
	"fmt"
	"strings"
)

var privileges = []string{
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"CREATE",
	"DROP",
	"REFERENCES",
	"INDEX",
	"ALTER",
	"CREATE TEMPORARY TABLES",
	"LOCK TABLES",
	"EXECUTE",
	"CREATE VIEW",
	"SHOW VIEW",
	"CREATE ROUTINE",
	"ALTER ROUTINE",
	"EVENT",
	"TRIGGER",
}

func (s *Shim) CreateUserGrant(name, password, database string) error {
	query := []string{
		fmt.Sprintf("CREATE USER '%v'@'%%' IDENTIFIED BY '%v';", name, password),
		fmt.Sprintf("GRANT %v ON %v.* TO '%v'@'%%';", strings.Join(privileges, ","), database, name),
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
