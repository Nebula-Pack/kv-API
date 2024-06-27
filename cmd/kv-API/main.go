package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Nebula-Pack/kv-API/internal/config"
	"github.com/Nebula-Pack/kv-API/internal/handlers"
	"github.com/Nebula-Pack/kv-API/internal/middleware"
)

func main() {
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

	// Admin routes
	http.Handle("/admin/overwrite", middleware.AdminAuth(http.HandlerFunc(handlers.AdminOverwriteHandler(db))))
	http.HandleFunc("/admin/login", handlers.AdminLoginHandler())

	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
