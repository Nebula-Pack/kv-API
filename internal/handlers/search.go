package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

type SearchResult struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func SearchHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")
		if query == "" {
			http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
			return
		}

		log.Printf("Received search query: %s\n", query)
		results, err := searchInDB(db, query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Search results: %+v\n", results)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(results)
	}
}

func searchInDB(db *sql.DB, query string) ([]SearchResult, error) {
	rows, err := db.Query("SELECT key, value FROM kvstore")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []SearchResult
	lowerQuery := strings.ToLower(query)

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}

		lowerKey := strings.ToLower(key)
		lowerValue := strings.ToLower(value)

		// Check for substring matches
		if strings.Contains(lowerKey, lowerQuery) || strings.Contains(lowerValue, lowerQuery) {
			results = append(results, SearchResult{Key: key, Value: value})
			log.Printf("Substring match found: key: %s, value: %s\n", key, value)
			continue
		}

		// If no substring match, use fuzzy matching with a lower threshold
		keyMatch := fuzzy.Ratio(lowerQuery, lowerKey)
		valueMatch := fuzzy.Ratio(lowerQuery, lowerValue)

		log.Printf("Checking key: %s, value: %s, keyMatch: %d, valueMatch: %d\n", key, value, keyMatch, valueMatch)

		if keyMatch > 50 || valueMatch > 50 { // Lowered threshold to 50
			results = append(results, SearchResult{Key: key, Value: value})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
