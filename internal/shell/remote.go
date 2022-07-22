package shell

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type Runner interface {
	Start(cmd string) error
	Wait() error
	StdoutPipe() (io.Reader, error)
	StderrPipe() (io.Reader, error)
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

func (s RemoteShell) Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error {
	if in == nil {
		in = new(bytes.Buffer)
	}

	cmd = fmt.Sprintf("sh -c \"%v\"", cmd)

	stdout, err := s.runner.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	stdoutBuffer := bufio.NewReader(stdout)

	stderr, err := s.runner.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to create stderr pipe: %w", err)
	}
	stderrBuffer := bufio.NewReader(stderr)

	stdin, err := s.runner.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdin pipe: %w", err)
	}
	stdinBuffer := bufio.NewWriter(stdin)

	if err := s.runner.Start(cmd); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := io.Copy(out, stdoutBuffer)
		return err
	})

	if combinedOutput {
		g.Go(func() error {
			_, err := io.Copy(out, stderrBuffer)
			return err
		})
	}

	g.Go(func() error {
		_, err := io.Copy(stdinBuffer, in)
		stdin.Close()
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}
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
