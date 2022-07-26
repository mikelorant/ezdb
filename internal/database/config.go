package database

import (
	"crypto/tls"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

const (
	DBMaxAllowedPacket = 0
	DBTLSConfig        = "preferred"
)

func NewMySQLConfig(cfg Config) *mysql.Config {
	if cfg.Port == 0 {
		cfg.Port = 3306
	}

	dbcfg := mysql.NewConfig()
	dbcfg.Net = cfg.Host
	dbcfg.Addr = fmt.Sprintf("%v:%v", cfg.Host, cfg.Port)
	dbcfg.User = cfg.User
	dbcfg.Passwd = cfg.Password
	dbcfg.DBName = cfg.Name
	dbcfg.TLSConfig = DBTLSConfig
	dbcfg.MaxAllowedPacket = DBMaxAllowedPacket

	return dbcfg
}

func NewPostgresConfig(cfg Config) *pgx.ConnConfig {
	if cfg.Name == "" {
		cfg.Name = "postgres"
	}

	if cfg.Port == 0 {
		cfg.Port = 5432
	}

	fallbackcfg := pgconn.FallbackConfig{
		Host:      cfg.Host,
		Port:      uint16(cfg.Port),
		TLSConfig: nil,
	}

	var dbcfg pgx.ConnConfig
	dbcfg.Host = cfg.Host
	dbcfg.Database = cfg.Name
	dbcfg.Port = uint16(cfg.Port)
	dbcfg.User = cfg.User
	dbcfg.Password = cfg.Password
	dbcfg.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	dbcfg.Fallbacks = []*pgconn.FallbackConfig{
		&fallbackcfg,
	}

	return &dbcfg
}
