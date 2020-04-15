package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: <progname> <hash1>")
	}

	hash1 := os.Args[1]

	db, err := sql.Open("sqlite3", "./db.sqlite3")
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	row := db.QueryRow("SELECT signature FROM signatures WHERE hash1 = ?", hash1)
	var signature string
	if err := row.Scan(&signature); err != nil {
		log.Fatal(err)
	}
	fmt.Print(signature)
}
