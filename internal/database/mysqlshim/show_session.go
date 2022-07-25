package mysqlshim

import (
	"fmt"
)

const (
	QueryShowSession = "SELECT id, user, host, db, command, time, state, info FROM information_schema.processlist;"
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
