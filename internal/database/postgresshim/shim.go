package postgresshim

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"github.com/imdario/mergo"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

type Shim struct {
	DB  *sql.DB
	cfg *pgx.ConnConfig
}

const (
	EmptyPostgresURL = "postgresql://"
)

func New(cfg *pgx.ConnConfig, dialFunc func(ctx context.Context, address string) (net.Conn, error)) (*Shim, error) {
	parsedCfg, err := pgx.ParseConfig(EmptyPostgresURL)
	if err != nil {
		return nil, fmt.Errorf("unbable to parse config: %w", err)
	}

	if err := mergo.Merge(parsedCfg, cfg, mergo.WithOverride); err != nil {
		return nil, fmt.Errorf("unable to merge config: %w", err)
	}
	parsedCfg.DialFunc = dialFuncWithNet(dialFunc)

	con := stdlib.GetConnector(*parsedCfg)

	return &Shim{
		DB:  sql.OpenDB(con),
		cfg: cfg,
	}, nil
}

func dialFuncWithNet(dialFunc func(ctx context.Context, address string) (net.Conn, error)) func(ctx context.Context, network, address string) (net.Conn, error) {
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return dialFunc(ctx, address)
	}
}
