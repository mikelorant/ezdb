package postgresshim

import (
	"fmt"
	"strings"
)

var (
	PostgresDumpCommand = "pg_dump"
	PostgresDumpOptions = []string{
		"--no-password",
	}
)

func (s *Shim) BackupCommand(hidden bool) string {
	var cmd []string

	if s.cfg.Password != "" {
		if hidden {
			cmd = append(cmd, fmt.Sprintf("PGPASSWORD='%v'", strings.Repeat("*", len(s.cfg.Password))))
		} else {
			cmd = append(cmd, fmt.Sprintf("PGPASSWORD='%v'", s.cfg.Password))
		}
	}
	cmd = append(cmd, PostgresDumpCommand)

	if s.cfg.User != "" {
		cmd = append(cmd, fmt.Sprintf("--username=%v", s.cfg.User))
	}

	cmd = append(cmd, fmt.Sprintf("--host=%v", s.cfg.Host))
	cmd = append(cmd, fmt.Sprintf("--port=%v", s.cfg.Port))
	cmd = append(cmd, fmt.Sprintf("--dbname=%v", s.cfg.Database))

	cmd = append(cmd, PostgresDumpOptions...)

	return strings.Join(cmd, " ")
}
