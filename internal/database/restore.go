package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"strings"

	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
)

func (cl *Client) Restore(name, filename string, storer Storer, verbose bool) error {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Restoring..."
	bar := progress.New(-1, desc, verbose)

	done := make(chan bool)

	r, w := io.Pipe()
	storer.Retrieve(w, filename, done)
	rb := bufio.NewReader(r)

	scanner := bufio.NewScanner(rb)

	gz, err := rb.Peek(2)
	if gz[0] == 31 && gz[1] == 139 {
		rs := compress.NewGzipDecompressor(rb)
		scanner = bufio.NewScanner(rs)
	}

	buf := make([]byte, 0, 50*1024*1024) // Create a 50MB buffer
	scanner.Buffer(buf, 10*1024*1024)    // Scan lines of up to 10MB

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %w", err)
	}
	var query string
	var rows int64
	for scanner.Scan() {
		bar.Add(len(scanner.Bytes()))

		// Trim leading and trailing spaces
		text := strings.Trim(scanner.Text(), " ")

		// Skip line if empty or starts with a comment
		if text == "" || strings.HasPrefix(text, "--") {
			continue
		}

		// If we have no existing query set it to text
		// If we have a query join it with a space separator
		if query == "" {
			query = text
		} else {
			query = fmt.Sprintf("%v %v", query, text)
		}

		// If the line has a suffix of ";" we execute the query
		if strings.HasSuffix(text, ";") {
			res, err := tx.Exec(sanitise(query))
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("transaction failed: %w", err)
			}
			rowsAffected, _ := res.RowsAffected()
			rows += rowsAffected
			// Clear the query
			query = ""
		}

		if err := scanner.Err(); err != nil {
			tx.Rollback()
			return fmt.Errorf("unknown scanner error: %w", err)
		}
	}
	tx.Commit()

	<-done
	bar.Finish()

	return nil
}

// TODO: make generic
func sanitise(str string) string {
	substr := "DEFINER=`admin`@`%`"
	newstr := "DEFINER=`infra01`@`%`"

	if strings.Contains(str, substr) {
		return strings.Replace(str, substr, newstr, 1)
	}

	return str
}
