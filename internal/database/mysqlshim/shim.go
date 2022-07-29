package mysqlshim

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
)

type Shim struct {
	DB  *sql.DB
	cfg *mysql.Config
}

func New(cfg *mysql.Config, dialFunc func(ctx context.Context, address string) (net.Conn, error)) (*Shim, error) {
	mysql.RegisterDialContext(cfg.Net, dialFunc)

	con, err := mysql.NewConnector(cfg)
	if err != nil {
		return &Shim{}, fmt.Errorf("unable to create new connector: %w", err)
	}

	return &Shim{
		DB:  sql.OpenDB(con),
		cfg: cfg,
	}, nil
}

func (s *Shim) query(query string) ([][]string, error) {
	var out [][]string

	rows, err := s.DB.Query(query)
	if err != nil {
		return out, fmt.Errorf("unable to query database: %w", err)
	}
	defer rows.Close()
	log.Printf("Executed query: %s\n", query)

	out, err = toStringRows(rows)
	if err != nil {
		return out, fmt.Errorf("unable to output rows: %w", err)
	}

	return out, nil
}

func (s *Shim) queryRow(q string, v any) error {
	row := s.DB.QueryRow(q)
	if err := row.Scan(v); err != nil {
		return fmt.Errorf("unable to query row: %w", err)
	}

	return nil
}

func (s *Shim) exec(q string) error {
	tx, err := s.DB.Begin()
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

func toStringRows(rows *sql.Rows) ([][]string, error) {
	defer rows.Close()

	var out [][]string

	cols, err := rows.Columns()
	if err != nil {
		return out, fmt.Errorf("unable to get columns: %w", err)
	}
	out = append(out, cols)

	row := make([][]byte, len(cols))
	rowPtr := make([]any, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}

	for rows.Next() {
		err = rows.Scan(rowPtr...)
		if err != nil {
			return out, fmt.Errorf("unable to scan rows: %w", err)
		}

		out = append(out, toString(row))
	}

	return out, nil
}

func toString(b [][]byte) []string {
	s := make([]string, len(b))
	for i, v := range b {
		s[i] = string(v)
	}
	return s
}
