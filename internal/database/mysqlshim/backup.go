package mysqlshim

import (
	"fmt"
	"strings"
)

var (
	MySQLDumpCommand = "mysqldump"
	MySQLDumpOptions = []string{
		"--compress",
		"--routines",
		"--lock-tables=false",
		"--net_buffer_length=16384",
		"--ssl=true",
	}
)

func (s *Shim) BackupCommand(hidden bool) string {
	var cmd []string

	cmd = append(cmd, MySQLDumpCommand)

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

	cmd = append(cmd, MySQLDumpOptions...)

	cmd = append(cmd, s.cfg.DBName)

	return strings.Join(cmd, " ")
}
