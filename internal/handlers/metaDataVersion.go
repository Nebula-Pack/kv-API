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

func MetadataVersionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		version := r.URL.Query().Get("version")
		if key == "" || version == "" {
			http.Error(w, "key and version are required", http.StatusBadRequest)
			return
		}

		// Check if the key exists in the database
		exists, err := models.KeyExists(db, key)
		if err != nil {
			http.Error(w, fmt.Sprintf("error checking key existence: %v", err), http.StatusInternalServerError)
			return
		}
		if !exists {
			http.Error(w, "key not found", http.StatusNotFound)
			return
		}

		// Try to open the version metadata file
		filePath := filepath.Join("./data/ver_metadata", key, fmt.Sprintf("%s.json", version))
		file, err := os.Open(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				// If the version metadata file does not exist, return the value from the database
				value, err := models.GetValue(db, key)
				if err != nil {
					http.Error(w, fmt.Sprintf("error retrieving value from database: %v", err), http.StatusInternalServerError)
					return
				}

				// Call the clone endpoint
				cloneResp, err := utils.CheckIsLua(value, version)
				if err != nil {
					http.Error(w, fmt.Sprintf("error calling clone endpoint: %v", err), http.StatusInternalServerError)
					return
				}

				// Check if isLua is false
				if !cloneResp.IsLua {
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte("invalid package"))
					return
				}

				// Save the response as JSON in ver_metadata directory
				err = saveVersionedMetadata(key, version, cloneResp)
				if err != nil {
					http.Error(w, fmt.Sprintf("error saving versioned metadata: %v", err), http.StatusInternalServerError)
					return
				}

				// Return the clone response
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(cloneResp)
				return
			}
			http.Error(w, fmt.Sprintf("error opening metadata file: %v", err), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		var metadata map[string]interface{}
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&metadata)
		if err != nil {
			http.Error(w, fmt.Sprintf("error decoding JSON payload: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(metadata)
		if err != nil {
			http.Error(w, fmt.Sprintf("error encoding JSON payload: %v", err), http.StatusInternalServerError)
		}
	}
}

func saveVersionedMetadata(key, version string, cloneResp utils.CloneResponse) error {
	versionDir := filepath.Join("./data/ver_metadata", key)
	err := os.MkdirAll(versionDir, os.ModePerm)
	if err != nil {
		return err
	}

	filePath := filepath.Join(versionDir, fmt.Sprintf("%s.json", version))
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	jsonData, err := json.MarshalIndent(cloneResp, "", "  ")
	if err != nil {
		return err
	}

	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}

	return nil
}
