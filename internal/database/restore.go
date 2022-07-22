package database

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/mikelorant/ezdb2/internal/compress"
	"github.com/mikelorant/ezdb2/internal/progress"
	"golang.org/x/sync/errgroup"
)

type RestoreReader interface {
	io.Reader
}

var (
	MySQLRestoreCommand = "mysql"
	MySQLRestoreOptions = []string{
		"--compress",
		"--ssl-mode=preferred",
		"--protocol=tcp",
	}
)

func (cl *Client) Restore(name, filename string, storer Storer, verbose bool) error {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Restoring..."
	bar := progress.New(-1, desc, verbose)

	done := make(chan bool)

	r, w := io.Pipe()
	storer.Retrieve(w, filename)
	rb := bufio.NewReader(r)

	scanner := bufio.NewScanner(rb)

	gz, err := rb.Peek(2)
	if err != nil {
		return fmt.Errorf("unable to check for compression: %w", err)
	}
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

func (cl *Client) RestoreCompat(name, filename string, storer Storer, shell Shell, verbose bool) ([]byte, error) {
	db := sql.OpenDB(cl.connector)
	defer db.Close()

	desc := "Restoring..."
	bar := progress.New(-1, desc, verbose)

	var buf bytes.Buffer

	// storer (w) -> (w) pipe (r) -> (r) teereader (w) -> (w) progressbar
	// 			     			                   (r) -> (r) buffer (r) -> (r) gzip (r) -> mysql

	r, w := io.Pipe()
	tr := io.TeeReader(r, bar)
	rb := bufio.NewReader(tr)

	g := new(errgroup.Group)

	g.Go(func() error {
		return storer.Retrieve(w, filename)
	})

	var rr RestoreReader
	rr = rb

	gz, err := rb.Peek(2)
	if err != nil {
		return nil, fmt.Errorf("unable to check for compression: %w", err)
	}
	if gz[0] == 31 && gz[1] == 139 {
		rs := compress.NewGzipDecompressor(rb)
		rr = rs
	}

	log.Println("Command:", cl.getRestoreCommand(true))
	if err := shell.Run(&buf, rr, cl.getRestoreCommand(false), true); err != nil {
		bar.Finish()
		out, _ := io.ReadAll(&buf)
		fmt.Print(string(out))
		return nil, fmt.Errorf("unable to run command: %w", err)
	}

	out, err := io.ReadAll(&buf)
	if err != nil {
		return nil, fmt.Errorf("unable to read output: %w", err)
	}

	if err := g.Wait(); err != nil {
		return nil, fmt.Errorf("store failure: %w", err)
	}

	bar.Finish()

	return out, nil
}

func (cl *Client) getRestoreCommand(hidden bool) string {
	var cmd []string

	cmd = append(cmd, MySQLRestoreCommand)

	if user := cl.config.User; user != "" {
		cmd = append(cmd, fmt.Sprintf("--user=%v", user))
	}

	if password := cl.config.Passwd; password != "" {
		if hidden {
			cmd = append(cmd, fmt.Sprintf("--password=%v", strings.Repeat("*", len(password))))
		} else {
			cmd = append(cmd, fmt.Sprintf("--password=%v", password))
		}
	}

	hostPort := strings.Split(cl.config.Addr, ":")
	cmd = append(cmd, fmt.Sprintf("--host=%v --port=%v", hostPort[0], hostPort[1]))

	cmd = append(cmd, fmt.Sprintf("--database=%v", cl.config.DBName))

	cmd = append(cmd, MySQLRestoreOptions...)

	return strings.Join(cmd, " ")
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
