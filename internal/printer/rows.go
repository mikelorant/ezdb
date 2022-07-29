package printer

import (
	"bytes"

	"github.com/rodaine/table"
)

func Rows(rows [][]string) string {
	var buffer bytes.Buffer

	header := []interface{}{}
	for _, c := range rows[0] {
		header = append(header, c)
	}

	tbl := table.New(header...)
	tbl.WithWriter(&buffer)
	tbl.SetRows(rows[1:])
	tbl.Print()

	return buffer.String()
}
