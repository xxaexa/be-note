package main

import (
	"database/sql"
	"fmt"
	"go-note/api"
	"go-note/db"
	"log"
)

func main() {

	db, err := db.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	server := api.NewAPIServer(fmt.Sprintf(":%s", "8080"), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}

}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DB: connect success!")
}
