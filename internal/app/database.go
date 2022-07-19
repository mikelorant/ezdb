package app

import (
	"context"
	"fmt"
	"net"

	"github.com/go-sql-driver/mysql"
	"github.com/mikelorant/ezdb2/internal/dbutil"
	"github.com/mikelorant/ezdb2/internal/structutil"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type Databases []Database

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Context  string `yaml:"context"`
	Tunnel   string `yaml:"tunnel"`
}

type DBOptions struct {
	name string
}

const (
	DBMaxAllowedPacket = 0
	DBTLSConfig        = "preferred"
)

func (d Database) String() string {
	out, _ := structutil.Sprint(d)
	return out
}

func WithDBName(name string) func(*DBOptions) {
	return func(d *DBOptions) {
		d.name = name
	}
}

func (a *App) GetDBClient(context string, dbOpts ...func(*DBOptions)) (*dbutil.Client, error) {
	var dbOptions DBOptions
	for _, o := range dbOpts {
		o(&dbOptions)
	}

	db, tun := a.Config.getContext(context)

	dbcfg := getDBConfig(db, tun, dbOptions.name)
	dial, err := getDialerFunc(tun)
	if err != nil {
		fmt.Errorf("unable to get dialer function: %w", err)
	}
	cl, err := dbutil.NewClient(dbcfg, dial)
	if err != nil {
		return nil, fmt.Errorf("unable to get new client: %w", err)
	}

	return cl, nil
}

func getDialerFunc(tun *Tunnel) (func(ctx context.Context, address string) (net.Conn, error), error) {
	dial := dialerFunc(&net.Dialer{})

	if isTunnel(tun) {
		tunnel, err := makeTunnel(tun)
		if err != nil {
			return nil, fmt.Errorf("unable to make tunnel: %w", err)
		}
		dial = dialerFunc(tunnel)
	}

	return dial, nil
}

func getDBConfig(db *Database, tun *Tunnel, name string) *mysql.Config {
	cfg := mysql.NewConfig()
	cfg.Net = db.Host
	cfg.Addr = fmt.Sprintf("%v:3306", db.Host)
	cfg.User = db.User
	cfg.Passwd = db.Password
	cfg.DBName = name
	cfg.TLSConfig = "preferred"
	cfg.MaxAllowedPacket = DBMaxAllowedPacket

	return cfg
}

func dialerFunc(dialer Dialer) func(ctx context.Context, address string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		return dialer.Dial("tcp", address)
	}
}
