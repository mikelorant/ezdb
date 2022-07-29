package postgresshim

import (
	"fmt"
)

const (
	QueryDBExists = "SELECT 1 AS result FROM pg_database WHERE datname='%v'"
)

func (s *Shim) IsDatabase(name string) bool {
	var dbname string

	query := fmt.Sprintf(QueryDBExists, name)

	s.queryRow(query, &dbname)

	return dbname == name
}
