package database

import (
	"database/sql"
	"fmt"
)

const (
	QueryDBExists = "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%v';"
)

func (cl *Client) IsDatabase(name string) bool {
	var dbname string

	db := sql.OpenDB(cl.connector)
	defer db.Close()

	query := fmt.Sprintf(QueryDBExists, name)

	row := db.QueryRow(query)
	row.Scan(&dbname)

	return dbname == name
}
