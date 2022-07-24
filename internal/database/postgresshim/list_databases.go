package postgresshim

import (
	"fmt"
)

const (
	QueryShowDatabases = "SELECT datname FROM pg_database;"
)

func (s *Shim) ListDatabases() ([]string, error) {
	var list []string

	var res string
	rows, err := s.DB.Query(QueryShowDatabases)
	if err != nil {
		return nil, fmt.Errorf("unable to get database list: %w", err)
	}
	for rows.Next() {
		rows.Scan(&res)
		list = append(list, res)
	}

	return list, nil
}
