package postgresshim

var (
	PostgresDumpCommand = "mysqldump"
	PostgresDumpOptions = []string{}
)

func (s *Shim) BackupCommand(hidden bool) string {
	return ""
}
