package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func SaveMetadataHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload map[string]interface{}
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}

		key, ok := payload["key"].(string)
		if !ok || key == "" {
			http.Error(w, "Key is missing or invalid", http.StatusBadRequest)
			return
		}

		// Create the metadata directory if it doesn't exist
		metadataDir := "./data/metadata"
		err = os.MkdirAll(metadataDir, os.ModePerm)
		if err != nil {
			http.Error(w, "Failed to create metadata directory", http.StatusInternalServerError)
			return
		}

		// Save the payload to a file
		filePath := filepath.Join(metadataDir, fmt.Sprintf("%s.json", key))
		file, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Failed to create metadata file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		jsonData, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			http.Error(w, "Failed to serialize JSON payload", http.StatusInternalServerError)
			return
		}

		_, err = file.Write(jsonData)
		if err != nil {
			http.Error(w, "Failed to write to metadata file", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Metadata saved successfully")
	}
}
