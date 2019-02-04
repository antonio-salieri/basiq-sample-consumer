package main

import (
	"fmt"
	"log"
	"os"

	"github.com/antonio-salieri/basiq-sample-consumer/client"
	"github.com/antonio-salieri/basiq-sample-consumer/cmd"
	"github.com/antonio-salieri/basiq-sample-consumer/service"
)

const (
	apiKeyEnvName = "BASIQ_API_KEY"
	// apiKey        = "a0b6cac3-cac1-463d-a3e2-c5519dec0fae:d1deaadf-4c66-4628-8050-63c724ade733"
)

var basiqClient client.Client

func main() {

	apiClient, err := client.NewBasiqClient(os.Getenv(apiKeyEnvName))
	if err != nil {
		log.Fatalf("Error creating Basiq client: %s", err)
	}

	transactionService := service.NewTransactionService(apiClient)

	if err := cmd.Execute(apiClient, transactionService); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
