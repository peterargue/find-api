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
	_ = godotenv.Load()

	apiURL := os.Getenv("FIND_API_URL")
	user := os.Getenv("FIND_API_USER")
	pass := os.Getenv("FIND_API_PASS")

	if apiURL == "" || user == "" || pass == "" {
		log.Fatal("FIND_API_URL, FIND_API_USER, and FIND_API_PASS must be set")
	}

	client := findapi.NewClient(user, pass, findapi.WithBaseURL(apiURL))
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
