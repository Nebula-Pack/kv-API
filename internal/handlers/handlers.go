// internal/handlers/handlers.go
package handlers

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/Nebula-Pack/kv-API/internal/models"
	"github.com/Nebula-Pack/kv-API/utils"
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
		repo := r.URL.Query().Get("value")
		if key == "" || repo == "" {
			http.Error(w, "key and repo are required", http.StatusBadRequest)
			return
		}

		isLua, err := utils.CheckIsLua(repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !isLua {
			http.Error(w, "key-value pair is not allowed", http.StatusForbidden)
			return
		}

		err = models.SetValue(db, key, repo)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "key already exists", http.StatusConflict)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		fmt.Fprintf(w, "key-value pair added successfully")
	}
}

func AdminOverwriteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		repo := r.URL.Query().Get("value")
		if key == "" || repo == "" {
			http.Error(w, "key and repo are required", http.StatusBadRequest)
			return
		}

		err := models.SetValue(db, key, repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "key-value pair overwritten successfully")
	}
}
