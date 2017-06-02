package main

import (
	"bufio"
	"database/sql"
	"flag"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var sqliteDatabase = flag.String("sqlite", "", "sqlite database to edit")
var _ = MainHook(func() error {
	if *sqliteDatabase == "" {
		return nil
	}

	db, err := sql.Open("sqlite3", *sqliteDatabase)
	if err != nil {
		return err
	}

	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return err
		}
		_, err = db.Exec(string(line))
		if err != nil {
			return err
		}
	}
})
