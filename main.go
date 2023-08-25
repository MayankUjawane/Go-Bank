package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MayankUjawane/gobank/controllers"
	"github.com/MayankUjawane/gobank/db"
	"github.com/MayankUjawane/gobank/token"
	"github.com/MayankUjawane/gobank/types"
	"github.com/joho/godotenv"
)

func main() {
	seed := flag.Bool("seed", false, "seed teh db")
	flag.Parse()

	// loading env file
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Some error occured in env file: %s", err)
	}

	// making connection with the database
	store, err := db.NewPostgresStore()
	if err != nil {
		log.Fatalf("error while making connection with db: %s", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// seed some accounts
	if *seed {
		fmt.Println("Seeding the database")
		seedAccounts(store)
	}

	tokenMaker := token.NewJWTMaker(os.Getenv("JWT_SECRET"))
	server := controllers.NewAPIServer(":3000", store, tokenMaker)
	server.SetupRouter(tokenMaker)
}

// for seeding just after runing the program, pass the seed in console
// ./bin/gobank --seed
func seedAccounts(s db.Storage) {
	seedAccount(s, "Mayank", "Ujawane", "Hello")
}

func seedAccount(store db.Storage, fname, lname, pw string) {
	account := types.NewAccount(fname, lname, pw)
	fmt.Println("new account number => ", account.Number)
	store.CreateAccount(account)
}
