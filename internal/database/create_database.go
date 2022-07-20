package database

import (
	"database/sql"
	"fmt"
)

func (cl *Client) CreateDatabase(name string) error {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	q := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %v;", name)

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	_, err = tx.Exec(q)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("transaction failed: %w", err)
	}
	tx.Commit()

	return nil
}
