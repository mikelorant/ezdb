package mysqlshim

import (
	"database/sql"
	"fmt"
)

const (
	QueryDBSize = "SELECT SUM(DATA_LENGTH) FROM INFORMATION_SCHEMA.TABLES WHERE TABLE_SCHEMA='%v';"
)

func (s *Shim) GetDatabaseSize(name string) (int64, error) {
	var size sql.NullInt64

	query := fmt.Sprintf(QueryDBSize, name)

	s.queryRow(query, &size)

	return size.Int64, nil
}
