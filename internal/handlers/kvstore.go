package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"kv-API/internal/models"
)

func GetHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key is missing", http.StatusBadRequest)
			return
		}

		value, err := models.GetValue(db, key)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "key not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		fmt.Fprintf(w, "%s", value)
	}
}

func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "key and value are required", http.StatusBadRequest)
			return
		}

		err := models.SetValue(db, key, value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "key-value pair added/updated successfully")
	}
}
