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

	rows, err := s.DB.Query(q)
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
