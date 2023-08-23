package main

import (
	"log"

	"github.com/MayankUjawane/gobank/controllers"
	"github.com/MayankUjawane/gobank/db"
)

func main() {
	// making connection with the database
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := controllers.NewAPIServer(":3000", store)
	server.Router()
}
