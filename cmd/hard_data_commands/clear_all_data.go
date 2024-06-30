package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Path to SQLite database file
	dbPath := "data/kvstore.db"

	// Clear SQLite database
	err := clearSQLiteDB(dbPath)
	if err != nil {
		log.Fatalf("Failed to clear SQLite database: %v", err)
	}
	fmt.Println("SQLite database cleared successfully.")

	// Path to metadata folder
	metadataFolderPath := "data/metadata"

	// Clear metadata folder
	err = clearMetadataFolder(metadataFolderPath)
	if err != nil {
		log.Fatalf("Failed to clear metadata folder: %v", err)
	}
	fmt.Println("Metadata folder cleared successfully.")
}

// Function to clear SQLite database
func clearSQLiteDB(dbPath string) error {
	err := os.Remove(dbPath)
	if err != nil {
		return fmt.Errorf("error removing SQLite database: %w", err)
	}
	return nil
}

// Function to clear metadata folder
func clearMetadataFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path)
			if err != nil {
				return fmt.Errorf("error removing file %s: %w", path, err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking through metadata folder: %w", err)
	}
	return nil
}
