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

	for _, q := range query {
		if err := s.exec(q); err != nil {
			return fmt.Errorf("unable to exec query: %w", err)
		}
	}

	return nil
}
