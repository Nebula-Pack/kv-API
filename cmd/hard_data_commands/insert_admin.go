package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./data/kvstore.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO users (username, email, password, is_admin) VALUES (?, ?, ?, 1)", "KeaganGilmore", "keagangilmore@gmail.com", "#Keagan4200")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initial admin user created")
}
