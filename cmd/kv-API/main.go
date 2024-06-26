package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Nebula-Pack/kv-API/internal/config"
	"github.com/Nebula-Pack/kv-API/internal/handlers"
)

func main() {
	// Initialize the database
	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Set up HTTP routes
	http.HandleFunc("/get", handlers.GetHandler(db))
	http.HandleFunc("/post", handlers.PostHandler(db))
	http.HandleFunc("/save-metadata", handlers.SaveMetadataHandler())
	http.HandleFunc("/keys", handlers.GetAllKeysHandler(db))
	http.HandleFunc("/metadata", handlers.GetMetadataHandler()) // New endpoint

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
