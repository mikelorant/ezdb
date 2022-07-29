package postgresshim

import (
	"fmt"
)

const (
	QueryShowSession = "SELECT pid ,datname ,usename ,application_name ,client_hostname ,client_port ,backend_start ,query_start ,query ,state FROM pg_stat_activity"
)

func (s *Shim) ShowSession() ([][]string, error) {
	out, err := s.query(QueryShowSession)
	if err != nil {
		return nil, fmt.Errorf("unable to get sessions: %w", err)
	}

	return out, nil
}
