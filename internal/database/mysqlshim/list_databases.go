package mysqlshim

import (
	"fmt"
)

const (
	QueryShowDatabases = "SHOW DATABASES;"
)

func (s *Shim) ListDatabases() ([]string, error) {
	var list []string

	out, err := s.Query(QueryShowDatabases)
	if err != nil {
		return nil, fmt.Errorf("unable to get database list: %w", err)
	}

	for _, v := range out {
		list = append(list, v[0])
	}

	return list, nil
}
