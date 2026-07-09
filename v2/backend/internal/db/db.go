package db

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

func InitDB(dbPath, dbName, schemaFile string) *sql.DB {
	if err := os.MkdirAll(dbPath, os.ModePerm); err != nil {
		log.Fatalf("failed to create database directory: %v", err)
	}

	dbFile := filepath.Join(dbPath, dbName)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	pragmas := []string{
		`PRAGMA foreign_keys = ON;`,
		`PRAGMA journal_mode = WAL;`,
		`PRAGMA synchronous = NORMAL;`,
		`PRAGMA busy_timeout = 5000;`,
	}

	for _, p := range pragmas {
		if _, err := db.Exec(p); err != nil {
			log.Fatalf("failed to execute pragma: %v", err)
		}
	}

	if err := execSchema(db, schemaFile); err != nil {
		log.Fatalf("failed to execute schema (%s): %v", schemaFile, err)
	}

	log.Println("database initialized successfully")

	return db
}

func execSchema(db *sql.DB, schemaFile string) error {
	data, err := os.ReadFile(schemaFile)
	if err != nil {
		return err
	}

	schema := string(data)

	stmts := strings.Split(schema, ";")
	for _, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}

	return nil
}

func CloseDB(db *sql.DB) {
	if db == nil {
		return
	}

	if err := db.Close(); err != nil {
		log.Printf("error closing database: %v", err)
		return
	}

	log.Println("database closed")
}
