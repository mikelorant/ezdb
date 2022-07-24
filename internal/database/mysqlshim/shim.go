package mysqlshim

import (
	"context"
	"database/sql"
	"fmt"
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
