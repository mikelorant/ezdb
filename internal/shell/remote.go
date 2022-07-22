package shell

import (
	"bufio"
	"fmt"
	"io"
)

type Runner interface {
	Start(cmd string) error
	Wait() error
	StdoutPipe() (io.Reader, error)
}

type RemoteShell struct {
	runner Runner
}

func NewRemoteShell(run Runner) *RemoteShell {
	return &RemoteShell{
		runner: run,
	}
}

func (s RemoteShell) Run(out io.Writer, cmd string) error {
	cmd = fmt.Sprintf("sh -c \"%v\"", cmd)

	r, err := s.runner.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	rb := bufio.NewReader(r)

	if err := s.runner.Start(cmd); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	if _, err := io.Copy(out, rb); err != nil {
		return fmt.Errorf("unable to copy output: %w", err)
	}

	if err := s.runner.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
