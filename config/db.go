package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	// Create a new SQLite database file in the data folder
	db, err := sql.Open("sqlite3", "./data/kvstore.db")
	if err != nil {
		return nil, err
	}

	// Create the key-value table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS kvstore (
		"key" TEXT NOT NULL PRIMARY KEY,
		"value" TEXT
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	log.Println("Database initialized and table created")
	return db, nil
}
