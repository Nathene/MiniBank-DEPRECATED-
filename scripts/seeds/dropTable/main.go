package main

import (
	"log"

	"github.com/Nathene/MiniBank/internal"
)

func main() {
	store, err := internal.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	store.Db.Exec("drop table if exists account")
}
