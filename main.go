package main

import "log"

func main() {
	store, err := NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}
	store.CreateAccountTable()
	server := NewAPIServer(":3000", store)
	server.Run()
}
