package database

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/rodaine/table"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Storer interface {
	Store(data io.Reader, filename string, done chan bool, result chan string) error
	Retrieve(data io.WriteCloser, filename string, done chan bool) error
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
