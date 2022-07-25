package postgresshim

import (
	"fmt"
)

const (
	QueryShowSession = "SELECT pid ,datname ,usename ,application_name ,client_hostname ,client_port ,backend_start ,query_start ,query ,state FROM pg_stat_activity"
)

func (s *Shim) ShowSession() ([][]string, error) {
	rows, err := s.DB.Query(QueryShowSession)
	if err != nil {
		return nil, fmt.Errorf("unable to get database list: %w", err)
	}
	defer rows.Close()

	out, err := output(rows)
	if err != nil {
		return out, fmt.Errorf("unable to output rows: %w", err)
	}

	return out, nil
}
