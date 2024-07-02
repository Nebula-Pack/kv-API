package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Nebula-Pack/kv-API/internal/config"
	"github.com/Nebula-Pack/kv-API/internal/handlers"
	"github.com/Nebula-Pack/kv-API/internal/middleware"
)

var startDate time.Time

func main() {
	startDate = time.Now()

	db, err := config.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/get", handlers.GetHandler(db))
	http.HandleFunc("/post", handlers.PostHandler(db))
	http.HandleFunc("/save-metadata", handlers.SaveMetadataHandler())
	http.HandleFunc("/keys", handlers.GetAllKeysHandler(db))
	http.HandleFunc("/metadata", handlers.GetMetadataHandler())
	http.HandleFunc("/metadata-version", handlers.MetadataVersionHandler(db)) // Updated to pass db

	// Admin routes
	http.Handle("/admin/overwrite", middleware.AdminAuth(http.HandlerFunc(handlers.AdminOverwriteHandler(db))))
	http.Handle("/admin/delete", middleware.AdminAuth(http.HandlerFunc(handlers.AdminDeleteHandler(db)))) // New delete route
	http.HandleFunc("/admin/login", handlers.AdminLoginHandler())

	fmt.Printf("Server started on %s\n", startDate.Format(time.RFC3339))
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
