package mysqlshim

import (
	"fmt"
)

const (
	QueryDBExists = "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%v';"
)

func (s *Shim) IsDatabase(name string) bool {
	var dbname string

	query := fmt.Sprintf(QueryDBExists, name)

	row := s.DB.QueryRow(query)
	row.Scan(&dbname)

	return dbname == name
}
