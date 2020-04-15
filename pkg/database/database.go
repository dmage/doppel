package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Signature struct {
	Signature string
	Hash1     string // sha256
	Hash2     uint32 // hash2
	Hash3     int    // length
}

type ErrSignatureNotFound struct {
	Hash1 string
}

func (e ErrSignatureNotFound) Error() string {
	return fmt.Sprintf("signature %s not found", e.Hash1)
}

type SignatureIterator struct {
	rows *sql.Rows
}

func (si *SignatureIterator) Next() (*Signature, bool) {
	if !si.rows.Next() {
		return nil, false
	}
	var s Signature
	if err := si.rows.Scan(&s.Signature, &s.Hash1, &s.Hash2, &s.Hash3); err != nil {
		panic(err)
	}
	return &s, true
}

func (si *SignatureIterator) Close() error {
	return si.rows.Close()
}

type DB struct {
	db *sql.DB
}

func New(name string) (*DB, error) {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}

func NewDefault() (*DB, error) {
	return New("./db.sqlite3")
}

func (d *DB) Close() error {
	return d.db.Close()
}

func (d *DB) PutSignature(s *Signature) error {
	var one int
	row := d.db.QueryRow("SELECT 1 FROM signatures WHERE hash1 = ?", s.Hash1)
	switch err := row.Scan(&one); err {
	case sql.ErrNoRows:
		// continue
	case nil:
		return nil
	default:
		return fmt.Errorf("unable to check if signature %s is stored: %s", s.Hash1, err)
	}

	_, err := d.db.Exec(
		"INSERT INTO signatures (signature, hash1, hash2, hash3) VALUES (?, ?, ?, ?)",
		s.Signature, s.Hash1, s.Hash2, s.Hash3,
	)
	if err != nil {
		return fmt.Errorf("unable to insert signature %s: %s", s.Hash1, err)
	}

	return nil
}

func (d *DB) GetSignature(hash1 string) (*Signature, error) {
	row := d.db.QueryRow("SELECT signature, hash1, hash2, hash3 FROM signatures WHERE hash1 = ?", hash1)
	var s Signature
	switch err := row.Scan(&s.Signature, &s.Hash1, &s.Hash2, &s.Hash3); err {
	case nil:
		// continue
	case sql.ErrNoRows:
		return nil, ErrSignatureNotFound{
			Hash1: hash1,
		}
	default:
		return nil, fmt.Errorf("unable to get signature %s: %s", s.Hash1, err)
	}
	return &s, nil
}

func (d *DB) ListSignatures() (*SignatureIterator, error) {
	rows, err := d.db.Query("SELECT signature, hash1, hash2, hash3 FROM signatures")
	if err != nil {
		return nil, err
	}
	return &SignatureIterator{rows: rows}, nil
}

func (d *DB) PutDocument(id, hash1 string) error {
	var one int
	row := d.db.QueryRow("SELECT 1 FROM documents WHERE id = ?", id)
	switch err := row.Scan(&one); err {
	case sql.ErrNoRows:
		// continue
	case nil:
		return nil
	default:
		return fmt.Errorf("unable to check if document %s is stored: %s", id, err)
	}

	_, err := d.db.Exec(
		"INSERT INTO documents (id, hash1) VALUES (?, ?)",
		id, hash1,
	)
	if err != nil {
		return fmt.Errorf("unable to insert document %s: %s", id, err)
	}

	return nil
}
