package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/hudaprs/RESTfulAPIScotchIO/app/controllers"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error Loading .env file")
	}

	app := controllers.App{}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)

	app.RunServer()
}
