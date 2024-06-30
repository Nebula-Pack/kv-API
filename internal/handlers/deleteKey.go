package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// AdminDeleteHandler handles the deletion of a key and its metadata
func AdminDeleteHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key is required", http.StatusBadRequest)
			return
		}

		// Delete the key from the database
		err := deleteKey(db, key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Delete the corresponding metadata file
		err = deleteMetadata(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "key and metadata deleted successfully")
	}
}

// deleteKey deletes a key from the database
func deleteKey(db *sql.DB, key string) error {
	_, err := db.Exec("DELETE FROM kvstore WHERE key = ?", key)
	if err != nil {
		return fmt.Errorf("failed to delete key: %v", err)
	}
	return nil
}

// deleteMetadata deletes the metadata file for a given key
func deleteMetadata(key string) error {
	metadataDir := "./data/metadata"
	filePath := filepath.Join(metadataDir, fmt.Sprintf("%s.json", key))

	err := os.Remove(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file to delete, ignore the error
		}
		return fmt.Errorf("failed to delete metadata file: %v", err)
	}

	return nil
}
