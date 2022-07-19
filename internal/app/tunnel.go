package app

import (
	"fmt"
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/sshutil"
	"github.com/mikelorant/ezdb2/internal/structutil"
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
	out, _ := structutil.Sprint(t)
	return out
}

func isTunnel(tun *Tunnel) bool {
	if tun == nil {
		return false
	}
	return true
}

func makeTunnel(tun *Tunnel) (*ssh.Client, error) {
	var s sshutil.SSH
	if err := s.Connect(tun.Key, tun.Host, tun.User); err != nil {
		return nil, fmt.Errorf("unable to connect to %v: %w", tun.Host, err)
	}
	out, err := s.Command(CmdHostname)
	if err != nil {
		return nil, fmt.Errorf("unable to run command: %v: %w", CmdHostname, err)
	}

	log.Printf("Tunnel succesfully connected: %v\n", strings.TrimSpace(out))

	return s.Client(), nil
}
