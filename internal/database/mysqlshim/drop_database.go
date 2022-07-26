package mysqlshim

import (
	"fmt"
)

const (
	QueryDropDatabase = "DROP DATABASE IF EXISTS %v;"
)

func (s *Shim) DropDatabase(name string) error {
	q := fmt.Sprintf(QueryDropDatabase, name)

	if err := s.exec(q); err != nil {
		return fmt.Errorf("unable to drop database: %v: %w", name, err)
	}

	return nil
}
