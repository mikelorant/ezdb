package mysqlshim

import (
	"fmt"
)

func (s *Shim) Query(query string) ([][]string, error) {
	out, err := s.query(query)
	if err != nil {
		return out, fmt.Errorf("unable to query database: %w", err)
	}

	return out, nil
}
