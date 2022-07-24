package postgresshim

import (
	"fmt"
)

const (
	QueryDBExists = "SELECT 1 AS result FROM pg_database WHERE datname='%v'"
)

func (s *Shim) IsDatabase(name string) bool {
	var res int

	query := fmt.Sprintf(QueryDBExists, name)

	row := s.DB.QueryRow(query)
	row.Scan(&res)

	return res == 1
}
