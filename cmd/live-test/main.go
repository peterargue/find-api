package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	findapi "github.com/peterargue/find-api"
)

func main() {
	// Load .env if it exists; ignore error if file is absent.
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	user := os.Getenv("FINDAPI_USERNAME")
	pass := os.Getenv("FINDAPI_PASSWORD")

	if user == "" || pass == "" {
		log.Fatal("FINDAPI_USERNAME, and FINDAPI_PASSWORD must be set")
	}

	client := findapi.NewClient(user, pass)
	svc := client.Flow

	suites := []Suite{
		BlocksSuite(svc),
		FTSuite(svc),
		NFTSuite(svc),
		AccountsSuite(svc),
		TransactionSuite(svc),
		EvmSuite(svc),
		NodesSuite(svc),
		ContractsSuite(svc),
	}

	ctx := context.Background()
	runner := NewRunner(suites)
	runner.Run(ctx)
	runner.PrintSummary()
	os.Exit(runner.ExitCode())
}
