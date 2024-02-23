package main

import (
	"fmt"
	"log"
	"os"

	eversendSdk "github.com/cetric32/eversend_go_sdk"
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

	wallets, err := eversendApp.Wallets.List()

	fmt.Println(wallets, err)

	// transaction, err := eversendApp.Payouts.Transaction("BP1801706452633548")
	transaction, err := eversendApp.Payouts.Quotation("UGX", 200, "momo", "KE", "KES", "SOURCE")

	// beneficiaries, err := eversendApp.Crypto.AddressTransactions("")

	fmt.Println(transaction, err)
}
