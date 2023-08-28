package util

import (
	"fmt"

	"github.com/MayankUjawane/gobank/db"
	"github.com/MayankUjawane/gobank/types"
)

// for seeding just after runing the program, pass the seed in console
// ./bin/gobank --seed
func SeedAccounts(s db.Storage) {
	seedAccount(s, "Mayank", "Ujawane", "Hello")
}

func seedAccount(store db.Storage, fname, lname, pw string) {
	account := types.NewAccount(fname, lname, pw)
	fmt.Println("new account number => ", account.Number)
	store.CreateAccount(account)
}
