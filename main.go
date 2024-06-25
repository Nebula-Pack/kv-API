package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Create a new SQLite database file
	db, err := sql.Open("sqlite3", "./kvstore.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the key-value table if not exists
	createTableSQL := `CREATE TABLE IF NOT EXISTS kvstore (
        "key" TEXT NOT NULL PRIMARY KEY,
        "value" TEXT
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Set up HTTP server
	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key is missing", http.StatusBadRequest)
			return
		}

		var value string
		err = db.QueryRow("SELECT value FROM kvstore WHERE key = ?", key).Scan(&value)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "key not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		fmt.Fprintf(w, "%s", value)
	})

	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		key := "hump"
		value := "vrld/hump"

		// Insert or replace the key-value pair into the database
		_, err := db.Exec("INSERT OR REPLACE INTO kvstore (key, value) VALUES (?, ?)", key, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "key-value pair added/updated successfully")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
