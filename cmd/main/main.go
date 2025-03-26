package main

import (
	"fmt"
	"log"

	mainService "github.com/grozaqueen/julse/internal/apps/main_service"
	"github.com/joho/godotenv"
)

const configFile = ".env"

func main() {
	err := godotenv.Load(configFile)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	server, err := mainService.NewServer()
	if err != nil {
		log.Fatal(fmt.Errorf("error occured when creating server, %w", err))
	}

	if err = server.Run(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
