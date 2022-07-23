package app

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/mikelorant/ezdb2/internal/database/mysqlshim"
	"github.com/mikelorant/ezdb2/internal/structprinter"
)

type Databases []Database

type Database struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
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
	return structprinter.Sprint(d)
}

func WithDBName(name string) func(*DBOptions) {
	return func(d *DBOptions) {
		d.name = name
	}
}

func (a *App) GetDBClient(context string, dbOpts ...func(*DBOptions)) (*mysqlshim.Client, error) {
	var dbOptions DBOptions
	for _, o := range dbOpts {
		o(&dbOptions)
	}

	db, tun := a.Config.getContext(context)

	dbcfg := getDBConfig(db, tun, dbOptions.name)

	dial, err := getDialerFunc(tun)
	if err != nil {
		return nil, fmt.Errorf("unable to get dialer function: %w", err)
	}

	cl, err := mysqlshim.NewClient(dbcfg, dial)
	if err != nil {
		return nil, fmt.Errorf("unable to get new client: %w", err)
	}

	return cl, nil
}

func getDBConfig(db *Database, tun *Tunnel, name string) *mysql.Config {
	if db.Port == 0 {
		db.Port = 3306
	}

	cfg := mysql.NewConfig()
	cfg.Net = db.Host
	cfg.Addr = fmt.Sprintf("%v:%v", db.Host, db.Port)
	cfg.User = db.User
	cfg.Passwd = db.Password
	cfg.DBName = name
	cfg.TLSConfig = "preferred"
	cfg.MaxAllowedPacket = DBMaxAllowedPacket

	return cfg
}
