package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetMetadataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "Key is required", http.StatusBadRequest)
			return
		}

		// Construct file path
		filePath := filepath.Join("./data/metadata", fmt.Sprintf("%s.json", key))

		// Open the file
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "Metadata not found", http.StatusNotFound)
			return
		}
		defer file.Close()

		// Decode JSON from file
		var metadata map[string]interface{}
		err = json.NewDecoder(file).Decode(&metadata)
		if err != nil {
			http.Error(w, "Failed to read metadata", http.StatusInternalServerError)
			return
		}

		// Update metadata with latest GET request info
		metadata["temporal_semantics"] = map[string]interface{}{
			"latest-get-request": time.Now().Format(time.RFC3339),
		}

		// Save updated metadata back to file
		file, err = os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create metadata file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		jsonData, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {
			http.Error(w, "Failed to serialize JSON payload", http.StatusInternalServerError)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			http.Error(w, "Failed to write to metadata file", http.StatusInternalServerError)
			return
		}

		// Return JSON response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(metadata)
	}
}
