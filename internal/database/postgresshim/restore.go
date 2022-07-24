package postgresshim

var (
	PostgresRestoreCommand = "psql"
	PostgresRestoreOptions = []string{}
)

func (s *Shim) RestoreCommand(hidden bool) string {
	return ""
}
