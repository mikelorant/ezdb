package database

import (
	"bytes"
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
	List() ([]string, error)
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
