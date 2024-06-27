package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Nebula-Pack/kv-API/internal/models"
	"github.com/Nebula-Pack/kv-API/utils"
)

func PostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		repo := r.URL.Query().Get("value")
		if key == "" || repo == "" {
			http.Error(w, "key and repo are required", http.StatusBadRequest)
			return
		}

		cloneResp, err := utils.CheckIsLua(repo)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !cloneResp.IsLua {
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

		err = saveMetadata(key, cloneResp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "key-value pair added successfully")
	}
}

func saveMetadata(key string, cloneResp utils.CloneResponse) error {
	metadata := map[string]interface{}{
		"key":  key,
		"data": cloneResp,
	}

	metadataDir := "./data/metadata"
	err := os.MkdirAll(metadataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating metadata directory: %v", err)
	}

	filePath := filepath.Join(metadataDir, fmt.Sprintf("%s.json", key))
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating metadata file: %v", err)
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializing JSON payload: %v", err)
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to metadata file: %v", err)
	}

	return nil
}

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
