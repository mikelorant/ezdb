package tunnel

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type Tunnel struct {
	Client *ssh.Client
	config Config
}

type Config struct {
	Key  string `yaml:"key"`
	User string `yaml:"user"`
	Host string `yaml:"host"`
}

func (t *Tunnel) Connect(keyfile string, host string, user string) error {
	sign, err := signer(keyfile)
	if err != nil {
		return fmt.Errorf("unable to get signer: %w", err)
	}

	cfg := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(sign),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	cl, err := ssh.Dial("tcp", fmt.Sprintf("%v:22", host), cfg)
	if err != nil {
		return fmt.Errorf("unable to dial: %w", err)
	}
	t.Client = cl

	return nil
}

func (t *Tunnel) Command(cmd string) (string, error) {
	out, err := t.exec(cmd)
	if err != nil {
		return "", fmt.Errorf("unable to issue command: %v: %w", cmd, err)
	}

	return out.String(), nil
}

func (t *Tunnel) exec(cmd string) (*bytes.Buffer, error) {
	sess, err := t.Client.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Failed to create session: %w", err)
	}
	defer sess.Close()

	b, err := run(sess, cmd)
	if err != nil {
		return nil, fmt.Errorf("unable to exec: %v: %w", cmd, err)
	}

	return b, nil
}

func signer(keyfile string) (ssh.Signer, error) {
	key, err := os.ReadFile(keyfile)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %w", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %w", err)
	}

	return signer, nil
}

func run(sess *ssh.Session, cmd string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run(cmd); err != nil {
		return nil, fmt.Errorf("unable to run command: %v: %w", cmd, err)
	}

	return &b, nil
}
