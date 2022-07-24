package postgresshim

import (
	"fmt"
)

const (
	QueryDropDatabase = "DROP DATABASE %v;"
)

func (s *Shim) DropDatabase(name string) error {
	q := fmt.Sprintf(QueryDropDatabase, name)

	if _, err := s.DB.Exec(q); err != nil {
		return fmt.Errorf("unable to exec query: %w", err)
	}

	return nil
}
