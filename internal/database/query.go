package database

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"

	"github.com/rodaine/table"
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

func Format(rows [][]string) string {
	var buffer bytes.Buffer

	header := []interface{}{}
	for _, c := range rows[0] {
		header = append(header, c)
	}

	tbl := table.New(header...)
	tbl.WithWriter(&buffer)
	tbl.SetRows(rows[1:])
	tbl.Print()

	return buffer.String()
}

func output(rows *sql.Rows) ([][]string, error) {
	var out [][]string

	cols, err := rows.Columns()
	if err != nil {
		return out, fmt.Errorf("unable to get columns: %w", err)
	}
	out = append(out, cols)

	rawResult := make([][]byte, len(cols))

	dest := make([]interface{}, len(cols))
	for i := range rawResult {
		dest[i] = &rawResult[i]
	}

	for rows.Next() {
		result := []string{}
		err = rows.Scan(dest...)
		if err != nil {
			return out, fmt.Errorf("unable to scan rows: %w", err)
		}

		for _, v := range rawResult {
			if v == nil {
				result = append(result, "")
			} else {
				result = append(result, string(v))
			}
		}

		out = append(out, result)
	}

	return out, nil
}
