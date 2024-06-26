// internal/models/models.go
package models

import "database/sql"

func GetValue(db *sql.DB, key string) (string, error) {
	var value string
	err := db.QueryRow("SELECT value FROM kvstore WHERE key = ?", key).Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func KeyExists(db *sql.DB, key string) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM kvstore WHERE key = ?)", key).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func SetValue(db *sql.DB, key string, value string) error {
	// Check if the key already exists
	exists, err := KeyExists(db, key)
	if err != nil {
		return err
	}
	if exists {
		return sql.ErrNoRows // Custom error to indicate the key already exists
	}

	_, err = db.Exec("INSERT INTO kvstore (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		return err
	}
	return nil
}
