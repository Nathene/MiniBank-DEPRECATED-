package main

import (
	"log"

	"github.com/Nathene/MiniBank/internal"
	"github.com/Nathene/MiniBank/pkg/api"
)

func main() {
	store, err := internal.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}
	// store.db.Query("drop table account")
	server := api.NewAPIServer(":3000", store)
	server.Run()
}
