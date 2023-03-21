package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	conn := os.Getenv("POSTGRES_CONNECTION_STRING")
	prodRepo, err := NewProductRepository(conn)
	if err != nil {
		log.Fatal(err)
	}
	controller := NewController(prodRepo)
	if err := controller.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
