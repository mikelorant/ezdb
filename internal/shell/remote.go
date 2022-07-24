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

type RemoteShell struct {
	session Session
}

func NewRemoteShell(sess Session) *RemoteShell {
	return &RemoteShell{
		session: sess,
	}
}

func (s RemoteShell) Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error {
	if in == nil {
		in = new(bytes.Buffer)
	}

	cmd = fmt.Sprintf("sh -c \"%v\"", cmd)

	stdout, err := s.session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	stdoutBuffer := bufio.NewReader(stdout)

	stderr, err := s.session.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to create stderr pipe: %w", err)
	}
	stderrBuffer := bufio.NewReader(stderr)

	stdin, err := s.session.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdin pipe: %w", err)
	}
	stdinBuffer := bufio.NewWriter(stdin)

	if err := s.session.Start(cmd); err != nil {
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

	if err := s.session.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
