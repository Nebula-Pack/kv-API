package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/Nebula-Pack/kv-API/internal/models"
)

func GetAllKeysHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		keys, err := models.GetAllKeys(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(keys)
	}
}
