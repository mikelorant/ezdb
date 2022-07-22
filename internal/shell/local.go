package shell

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"golang.org/x/sync/errgroup"
)

type LocalShell struct{}

func NewLocalShell() *LocalShell {
	return &LocalShell{}
}

func (s LocalShell) Run(out io.Writer, in io.Reader, cmd string) error {
	if in == nil {
		in = new(bytes.Buffer)
	}

	c := exec.Command("sh", "-c", cmd)

	r, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	rb := bufio.NewReader(r)

	w, err := c.StdinPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdin pipe: %w", err)
	}
	wb := bufio.NewWriter(w)

	if err := c.Start(); err != nil {
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

	if err := c.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
