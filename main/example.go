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

	quotation, error := eversendApp.CreatePayoutQuotation("UGX", 5000, "momo", "KE", "KES", "DESTINATION")

	fmt.Println(quotation, error)

	token := quotation["token"].(string)

	transaction, error := eversendApp.CreateMomoPayout(token, "+254797224768", "cetric", "lihalakha", "KE")

	fmt.Println(transaction, error)
}
