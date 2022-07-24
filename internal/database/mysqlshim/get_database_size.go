package mysqlshim

import (
	"database/sql"
	"fmt"
)

const (
	QueryDBSize = "select SUM(DATA_LENGTH) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%v';"
)

func (s *Shim) GetDatabaseSize(name string) (int64, error) {
	var size sql.NullInt64

	query := fmt.Sprintf(QueryDBSize, name)

	row := s.DB.QueryRow(query)
	if err := row.Scan(&size); err != nil {
		return size.Int64, fmt.Errorf("unable to get database size: %w", err)
	}

	return size.Int64, nil
}
