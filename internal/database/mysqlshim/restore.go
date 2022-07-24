package mysqlshim

import (
	"fmt"
	"strings"
)

var (
	MySQLRestoreCommand = "mysql"
	MySQLRestoreOptions = []string{
		"--compress",
		// "--ssl-mode=preferred",
		"--protocol=tcp",
	}
)

func (s *Shim) RestoreCommand(hidden bool) string {
	var cmd []string

	cmd = append(cmd, MySQLRestoreCommand)

	if s.cfg.User != "" {
		cmd = append(cmd, fmt.Sprintf("--user=%v", s.cfg.User))
	}

	if s.cfg.Passwd != "" {
		if hidden {
			cmd = append(cmd, fmt.Sprintf("--password=%v", strings.Repeat("*", len(s.cfg.Passwd))))
		} else {
			cmd = append(cmd, fmt.Sprintf("--password=%v", s.cfg.Passwd))
		}
	}

	hostPort := strings.Split(s.cfg.Addr, ":")
	cmd = append(cmd, fmt.Sprintf("--host=%v --port=%v", hostPort[0], hostPort[1]))

	cmd = append(cmd, fmt.Sprintf("--database=%v", s.cfg.DBName))

	cmd = append(cmd, MySQLRestoreOptions...)

	return strings.Join(cmd, " ")
}
