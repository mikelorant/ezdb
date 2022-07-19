package database

import (
	"database/sql"
	"fmt"
	"log"
)

func (cl *Client) Query(query string) ([][]string, error) {
	var out [][]string

	db := sql.OpenDB(cl.connector)
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return out, fmt.Errorf("unable to query database: %w", err)
	}
	defer rows.Close()
	log.Printf("Executed query: %s\n", query)

	out, err = output(rows)
	if err != nil {
		return out, fmt.Errorf("unable to output rows: %w", err)
	}

	return out, nil
}
