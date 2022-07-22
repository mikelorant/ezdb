package shell

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

type LocalShell struct{}

func NewLocalShell() *LocalShell {
	return &LocalShell{}
}

func (s LocalShell) Run(out io.Writer, cmd string) error {
	c := exec.Command("sh", "-c", cmd)

	r, err := c.StdoutPipe()
	if err != nil {
		return fmt.Errorf("unable to create stdout pipe: %w", err)
	}
	rb := bufio.NewReader(r)

	if err := c.Start(); err != nil {
		return fmt.Errorf("unable to run command: %w", err)
	}

	if _, err := io.Copy(out, rb); err != nil {
		return fmt.Errorf("unable to copy output: %w", err)
	}

	if err := c.Wait(); err != nil {
		return fmt.Errorf("error running command: %w", err)
	}

	return nil
}
