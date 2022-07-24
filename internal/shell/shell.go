package shell

import "io"

type Runner interface {
	Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error
}

type Session interface {
	Start(cmd string) error
	Wait() error
	StdoutPipe() (io.Reader, error)
	StderrPipe() (io.Reader, error)
	StdinPipe() (io.WriteCloser, error)
}

type Shell struct {
	Config Config
	runner Runner
}

type Config struct {
	Session Session
}

func New(cfg Config) (*Shell, error) {
	if cfg.Session == nil {
		return &Shell{
			runner: NewLocalShell(),
		}, nil
	}

	return &Shell{
		runner: NewRemoteShell(cfg.Session),
	}, nil
}

func (s *Shell) Run(out io.Writer, in io.Reader, cmd string, combinedOutput bool) error {
	return s.runner.Run(out, in, cmd, combinedOutput)
}
