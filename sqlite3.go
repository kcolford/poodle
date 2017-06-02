package poodle

import (
	"bufio"
	"database/sql"
	"io"
)

// SqliteInterpreter runs a sqlite3 interpreter to be interacted with.
func SqliteInterpreter(dbname string, stdin io.Reader) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		return err
	}

	r := bufio.NewReader(stdin)
	for {
		line, err := r.ReadString(';')
		if err != nil {
			return err
		}
		_, err = db.Exec(string(line))
		if err != nil {
			return err
		}
	}
	return nil
}
