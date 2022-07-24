package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/structprinter"
	"github.com/mikelorant/ezdb2/internal/tunnel"
	"golang.org/x/crypto/ssh"
)

type Tunnels []Tunnel

type Tunnel struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
	User string `yaml:"user"`
	Host string `yaml:"host"`
}

const (
	CmdHostname = "/bin/hostname"
)

func (t Tunnel) String() string {
	return structprinter.Sprint(t)
}

func isTunnel(tun *Tunnel) bool {
	return tun != nil
}

func makeTunnel(tun *Tunnel) (*ssh.Client, error) {
	var t tunnel.Tunnel
	if err := t.Connect(tun.Key, tun.Host, tun.User); err != nil {
		return nil, fmt.Errorf("unable to connect to %v: %w", tun.Host, err)
	}
	out, err := t.Command(CmdHostname)
	if err != nil {
		return nil, fmt.Errorf("unable to run command: %v: %w", CmdHostname, err)
	}

	log.Printf("Tunnel succesfully connected: %v\n", strings.TrimSpace(out))

	return t.Client, nil
}

func getTunnelSession(tun *Tunnel) (*ssh.Session, error) {
	tunnel, err := makeTunnel(tun)
	if err != nil {
		return nil, fmt.Errorf("unable to make tunnel: %w", err)
	}

	sess, err := tunnel.NewSession()
	if err != nil {
		return nil, fmt.Errorf("unable to create a tunnel session: %w", err)
	}

	return sess, nil
}
