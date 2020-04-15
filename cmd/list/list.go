package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatal("usage: <progname>")
	}

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM signatures")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(id)
	}
}
