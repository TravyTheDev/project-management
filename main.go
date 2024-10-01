package main

import (
	"log"
	"project-management/api"
	"project-management/db"
)

func main() {
	db, err := db.NewSqlStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := api.NewAPIServer(":8000", db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
