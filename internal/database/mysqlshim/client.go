package mysqlshim

import (
	"database/sql/driver"
	"fmt"
	"io"

	"github.com/go-sql-driver/mysql"
)

type Storer interface {
	Store(data io.Reader, filename string) (string, error)
	Retrieve(data io.WriteCloser, filename string) error
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
