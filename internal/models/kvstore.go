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

func SetValue(db *sql.DB, key string, value string) error {
	_, err := db.Exec("INSERT OR REPLACE INTO kvstore (key, value) VALUES (?, ?)", key, value)
	if err != nil {
		return err
	}
	return nil
}
