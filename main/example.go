package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cetric32/eversend_go_sdk/eversendSdk"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", error(err))
	}

	clientId := os.Getenv("EVERSEND_CLIENT_ID")
	clientSecret := os.Getenv("EVERSEND_CLIENT_SECRET")

	eversendApp := eversendSdk.NewEversendApp(clientId, clientSecret)

	// transaction, err := eversendApp.GetTransaction("BP1801708338760630")

	beneficiaries, err := eversendApp.Crypto.GetAddresses()

	fmt.Println(beneficiaries, err)
}
