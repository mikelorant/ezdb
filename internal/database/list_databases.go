package database

import (
	"database/sql"
	"fmt"
)

const (
	QueryShowDatabases = "SHOW DATABASES;"
)

func (cl *Client) ListDatabases() ([]string, error) {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	var list []string

	var res string
	rows, err := db.Query(QueryShowDatabases)
	if err != nil {
		return nil, fmt.Errorf("unable to get database list: %w", err)
	}
	for rows.Next() {
		rows.Scan(&res)
		list = append(list, res)
	}

	return list, nil
}
