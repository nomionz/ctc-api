package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nomionz/ctc-api/api"
	"github.com/nomionz/ctc-api/repositories"
)

func main() {
	godotenv.Load()
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	prodRepo, err := repositories.NewProductRepository(conn)
	if err != nil {
		log.Fatal(err)
	}
	controller := api.NewController(prodRepo)
	if err := controller.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
