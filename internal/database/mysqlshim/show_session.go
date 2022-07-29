package mysqlshim

import (
	"fmt"
)

const (
	QueryShowSession = "SELECT id, user, host, db, command, time, state, info FROM information_schema.processlist;"
)

func (s *Shim) ShowSession() ([][]string, error) {
	out, err := s.query(QueryShowSession)
	if err != nil {
		return nil, fmt.Errorf("unable to get sessions: %w", err)
	}

	return out, nil
}
