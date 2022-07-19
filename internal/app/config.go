package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Databases Databases
	Tunnels   Tunnels
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

	if err := yaml.NewDecoder(file).Decode(c); err != nil {
		return fmt.Errorf("unable to decode config file: %v: %w", ConfigFile, err)
	}

	return nil
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

func (c *Config) getTunnel(tunnel string) *Tunnel {
	for _, v := range c.Tunnels {
		if v.Name == tunnel {
			return &v
		}
	}

	return nil
}
