package shell

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"

	"golang.org/x/sync/errgroup"
)

type LocalShell struct{}

func NewLocalShell() *LocalShell {
	return &LocalShell{}
}

func (s LocalShell) Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error {
	if in == nil {
		in = new(bytes.Buffer)
	}

	c := exec.Command("sh", "-c", cmd)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}

	stderr, err := c.StderrPipe()
	if err != nil {
		return fmt.Errorf("unable to create stderr pipe: %w", err)
	}

	stdin, err := c.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdin pipe: %w", err)
	}

	if err := c.Start(); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	g := new(errgroup.Group)

	g.Go(func() error {
		_, err := io.Copy(out, stdout)
		return err
	})

	if combinedOutput {
		g.Go(func() error {
			_, err := io.Copy(out, stderr)
			return err
		})
	} else {
		g.Go(func() error {
			_, err := io.Copy(os.Stderr, stderr)
			return err
		})
	}

	g.Go(func() error {
		_, err := io.Copy(stdin, in)
		stdin.Close()
		if errors.Is(err, syscall.EPIPE) {
			return nil
		}
		return err
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("unable to copy stdout/stdin: %w", err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
