package mysqlshim

import (
	"fmt"
)

const (
	QueryShowVariable    = "SHOW VARIABLES LIKE '%v';"
	QueryShowVariableAll = "SHOW VARIABLES;"
)

func (s *Shim) ShowVariable(variable string) ([][]string, error) {
	q := QueryShowVariableAll
	if variable != "" {
		q = fmt.Sprintf(QueryShowVariable, variable)
	}

	out, err := s.query(q)
	if err != nil {
		return nil, fmt.Errorf("unable to get variable: %w", err)
	}

	return out, nil
}
