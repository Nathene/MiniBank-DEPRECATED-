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

func seedAccount(store internal.Storage, fname, lname, user, email, pw string) *util.Account {
	acc, err := util.NewAccount(fname, lname, user, email, pw)
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
	seedAccount(s, "Nathan", "gg", "user", "nathan@email.com", "hehe")
	seedAccount(s, "Jessie", "Bessy", "user", "jess@email.com", "haha")
	seedAccount(s, "Bob", "J", "user", "bob@email.com", "hoho")
	seedAccount(s, "Fred", "P", "user", "fred@email.com", "huhu")
}
