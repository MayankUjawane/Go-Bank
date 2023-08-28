package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/MayankUjawane/gobank/controllers"
	"github.com/MayankUjawane/gobank/db"
	"github.com/MayankUjawane/gobank/token"
	"github.com/MayankUjawane/gobank/util"
	"github.com/joho/godotenv"
)

func main() {
	// storing seed value as false, will pass true through command prompt whenever required to seed the value
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

	// creating required table in database
	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	// seed some accounts
	if *seed {
		fmt.Println("Seeding the database")
		util.SeedAccounts(store)
	}

	tokenMaker := token.NewJWTMaker(os.Getenv("JWT_SECRET"))
	server := controllers.NewAPIServer(":3000", store, tokenMaker)
	server.SetupRouter(tokenMaker)
}
