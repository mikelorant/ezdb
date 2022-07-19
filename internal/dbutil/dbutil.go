package dbutil

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/rodaine/table"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Client struct {
	config    *mysql.Config
	connector driver.Connector
}

type Config struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

var privileges = []string{
	"SELECT",
	"INSERT",
	"UPDATE",
	"DELETE",
	"CREATE",
	"DROP",
	"REFERENCES",
	"INDEX",
	"ALTER",
	"CREATE TEMPORARY TABLES",
	"LOCK TABLES",
	"EXECUTE",
	"CREATE VIEW",
	"SHOW VIEW",
	"CREATE ROUTINE",
	"ALTER ROUTINE",
	"EVENT",
	"TRIGGER",
}

func NewClient(cfg *mysql.Config, dialer mysql.DialContextFunc) (*Client, error) {
	mysql.RegisterDialContext(cfg.Net, dialer)

	con, err := mysql.NewConnector(cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to create new connector: %w", err)
	}

	return &Client{
		config:    cfg,
		connector: con,
	}, nil
}

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

func (cl *Client) CreateUserGrant(name, password, database string) error {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	query := []string{
		fmt.Sprintf("CREATE USER '%v'@'%%' IDENTIFIED BY '%v';", name, password),
		fmt.Sprintf("GRANT %v ON %v.* TO '%v'@'%%';", strings.Join(privileges, ","), database, name),
	}

	tx, err := db.Begin()
	for _, q := range query {
		_, err = tx.Exec(q)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("transaction failed: %w", err)
		}
	}
	tx.Commit()

	return nil
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
			return out, fmt.Errorf("unable to scan rows: %v", err)
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
