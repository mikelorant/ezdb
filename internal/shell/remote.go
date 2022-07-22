package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Start(cmd string) error
	Wait() error
	StdoutPipe() (io.Reader, error)
	StdinPipe() (io.WriteCloser, error)
}

type RemoteShell struct {
	runner Runner
}

func NewRemoteShell(run Runner) *RemoteShell {
	return &RemoteShell{
		runner: run,
	}
}

func (s RemoteShell) Run(out io.Writer, in io.Reader, cmd string) error {
	if in == nil {
		in = new(bytes.Buffer)
	}

	cmd = fmt.Sprintf("sh -c \"%v\"", cmd)

	r, err := s.runner.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	rb := bufio.NewReader(r)

	w, err := s.runner.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdin pipe: %w", err)
	}
	wb := bufio.NewWriter(w)

	if err := s.runner.Start(cmd); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := io.Copy(out, rb)
		return err
	})

	g.Go(func() error {
		_, err := io.Copy(wb, in)
		return err
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("unable to copy stdout/stdin: %w", err)
	}

	if err := s.runner.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
