package sshutil

import (
	"bytes"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSH struct {
	client *ssh.Client
	config Config
}

type Config struct {
	Key  string `yaml:"key"`
	User string `yaml:"user"`
	Host string `yaml:"host"`
}

func (s *SSH) Connect(keyfile string, host string, user string) error {
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
	s.client = cl

	return nil
}

func (s *SSH) Client() *ssh.Client {
	return s.client
}

func (s *SSH) Command(cmd string) (string, error) {
	out, err := s.exec(cmd)
	if err != nil {
		return "", fmt.Errorf("unable to issue command: %v: %w", cmd, err)
	}

	return out.String(), nil
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

func (s *SSH) exec(cmd string) (*bytes.Buffer, error) {
	sess, err := s.client.NewSession()
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

func run(sess *ssh.Session, cmd string) (*bytes.Buffer, error) {
	var b bytes.Buffer
	sess.Stdout = &b
	if err := sess.Run(cmd); err != nil {
		return nil, fmt.Errorf("unable to run command: %v: %w", cmd, err)
	}

	return &b, nil
}
