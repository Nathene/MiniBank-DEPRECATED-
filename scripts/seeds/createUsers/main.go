package main

import (
	"fmt"
	"log"

	"github.com/Nathene/MiniBank/internal"
	"github.com/Nathene/MiniBank/pkg/util"
)

func main() {
	store, err := internal.NewPostGresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	seedAccounts(store)

}

func seedAccount(store internal.Storage, fname, lname, email, pw string) *util.Account {
	acc, err := util.NewAccount(fname, lname, email, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("new account => \nName: %s\nNumber: %d\n", acc.FirstName+" "+acc.LastName, acc.Number)

	return acc
}

func seedAccounts(s internal.Storage) {
	seedAccount(s, "Nathan", "gg", "nathan@email.com", "hehe")
	seedAccount(s, "Jessie", "Bessy", "jess@email.com", "haha")
	seedAccount(s, "Bob", "J", "bob@email.com", "hoho")
	seedAccount(s, "Fred", "P", "fred@email.com", "huhu")
}
