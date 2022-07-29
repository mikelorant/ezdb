package postgresshim

import (
	"fmt"
)

func (s *Shim) CreateUserGrant(name, password, database string) error {
	query := []string{
		fmt.Sprintf("CREATE USER %v WITH PASSWORD '%v';", name, password),
		fmt.Sprintf("GRANT ALL ON DATABASE %v TO %v;", database, name),
	}

	for _, q := range query {
		if err := s.exec(q); err != nil {
			return fmt.Errorf("unable to exec query: %w", err)
		}
	}

	return nil
}
