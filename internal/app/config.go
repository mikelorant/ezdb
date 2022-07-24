package app

import (
	"fmt"
	"io"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases Databases
	Tunnels   Tunnels
	Stores    Stores
}

const (
	ConfigFile = ".ezdb.yaml"
)

func newConfig() (Config, error) {
	var c Config
	if err := c.load(); err != nil {
		return c, fmt.Errorf("unable to load config: %w", err)
	}

	return c, nil
}

func (c *Config) load() error {
	file, err := os.Open(ConfigFile)
	if err != nil {
		return fmt.Errorf("unable to read config file: %v: %w", ConfigFile, err)
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("unable to read contents of config file: %w", err)
	}

	r := strings.NewReader(os.ExpandEnv(string(data)))

	if err := yaml.NewDecoder(r).Decode(c); err != nil {
		return fmt.Errorf("unable to decode config file: %v: %w", ConfigFile, err)
	}

	return nil
}

func (c *Config) getContexts() []string {
	var res []string
	for _, v := range c.Databases {
		res = append(res, v.Context)
	}
	return res
}

func (c *Config) getStores() []string {
	var res []string
	for _, v := range c.Stores {
		res = append(res, v.Name)
	}
	return res
}

func (c *Config) getContext(context string) (database *Database, tunnel *Tunnel) {
	db := c.getDatabase(context)
	if db.Tunnel == "" {
		return db, nil
	}

	return db, c.getTunnel(db.Tunnel)
}

func (c *Config) getDatabase(context string) *Database {
	for _, v := range c.Databases {
		if v.Context == context {
			return &v
		}
	}

	return nil
}

func (c *Config) getTunnel(name string) *Tunnel {
	for _, v := range c.Tunnels {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

func (c *Config) getStore(name string) *Store {
	for _, v := range c.Stores {
		if v.Name == name {
			return &v
		}
	}

	return nil
}
