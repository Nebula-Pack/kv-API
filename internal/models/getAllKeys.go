package models

import "database/sql"

func GetAllKeys(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT key, value FROM kvstore")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kvMap := make(map[string]string)
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		kvMap[key] = value
	}

	return kvMap, nil
}
