package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/dmage/doppel/pkg/database"
	"github.com/dmage/doppel/pkg/hash2"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage: <progname> <id>")
	}

	db, err := database.NewDefault()
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	id := os.Args[1]

	buf, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	hash1 := sha256.Sum256(buf)
	hash1hex := fmt.Sprintf("%02x", hash1)

	hsh2, err := hash2.Hash2(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	hash3, err := hash2.Hash3(bytes.NewReader(buf))
	if err != nil {
		log.Fatal(err)
	}

	log.Println(id)
	err = db.PutSignature(&database.Signature{
		Signature: string(buf),
		Hash1:     hash1hex,
		Hash2:     hsh2,
		Hash3:     hash3,
	})
	if err != nil {
		log.Fatal("failed to insert record:", err)
	}

	err = db.PutDocument(id, hash1hex)
	if err != nil {
		log.Fatal("failed to insert document:", err)
	}
}
