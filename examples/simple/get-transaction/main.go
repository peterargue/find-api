package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	findapi "github.com/peterargue/find-api"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	// Get credentials from environment variables
	username := os.Getenv("FINDAPI_USERNAME")
	password := os.Getenv("FINDAPI_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("FINDAPI_USERNAME and FINDAPI_PASSWORD environment variables must be set")
	}

	// Create client
	client := findapi.NewClient(username, password)

	ctx := context.Background()

	// Get a specific transaction by ID
	transactionID := "b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562"

	fmt.Printf("Fetching transaction: %s\n\n", transactionID)

	txResp, err := client.Simple.GetTransaction().
		ID(transactionID).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get transaction: %v", err)
	}

	// Display transaction information
	if len(txResp.Transactions) == 0 {
		fmt.Println("No transaction found")
		return
	}

	tx := txResp.Transactions[0]

	fmt.Printf("Transaction Details:\n")
	fmt.Printf("  ID: %s\n", tx.ID)
	fmt.Printf("  Status: %s\n", tx.Status)
	fmt.Printf("  Block Height: %d\n", tx.BlockHeight)
	fmt.Printf("  Block ID: %s\n", tx.BlockID)
	fmt.Printf("  Timestamp: %s\n", tx.Timestamp)
	fmt.Println()

	fmt.Printf("Gas and Fees:\n")
	fmt.Printf("  Gas Limit: %d\n", tx.GasLimit)
	fmt.Printf("  Gas Used: %d\n", tx.GasUsed)
	fmt.Printf("  Fee: %.8f\n", tx.Fee)
	fmt.Println()

	fmt.Printf("Participants:\n")
	fmt.Printf("  Payer: %s\n", tx.Payer)
	fmt.Printf("  Proposer: %s\n", tx.Proposer)
	fmt.Printf("  Proposer Index: %d\n", tx.ProposerIndex)
	fmt.Printf("  Proposer Sequence Number: %d\n", tx.ProposerSequenceNumber)
	if len(tx.Authorizers) > 0 {
		fmt.Printf("  Authorizers:\n")
		for i, auth := range tx.Authorizers {
			fmt.Printf("    %d. %s\n", i+1, auth)
		}
	}
	fmt.Println()

	// Display error if present
	if tx.Error != "" {
		fmt.Printf("Error:\n")
		fmt.Printf("  Message: %s\n", tx.Error)
		fmt.Printf("  Code: %s\n", tx.ErrorCode)
		fmt.Println()
	}

	// Display events
	if len(tx.Events) > 0 {
		fmt.Printf("Events (%d total):\n", len(tx.Events))
		for i, event := range tx.Events {
			fmt.Printf("  Event #%d:\n", i+1)
			fmt.Printf("    Name: %s\n", event.Name)
			fmt.Printf("    Index: %d\n", event.EventIndex)
			if len(event.Fields) > 0 {
				fieldsJSON, err := json.MarshalIndent(event.Fields, "      ", "  ")
				if err != nil {
					fmt.Printf("      (error formatting fields: %v)\n", err)
				} else {
					fmt.Printf("      Fields: %s\n", string(fieldsJSON))
				}
			}
			fmt.Println()

			// Only show first 3 events
			if i >= 2 && len(tx.Events) > 3 {
				fmt.Printf("  ... and %d more event(s)\n\n", len(tx.Events)-3)
				break
			}
		}
	}

	// Display transaction body if present
	if tx.TransactionBody != nil && tx.TransactionBody.Body != "" {
		fmt.Printf("Transaction Body (Cadence Script):\n")
		fmt.Printf("---\n%s\n---\n", tx.TransactionBody.Body)
	}

	// Display events aggregate if present
	if len(tx.EventsAggregate) > 0 {
		fmt.Printf("Events Aggregate:\n")
		aggJSON, err := json.MarshalIndent(tx.EventsAggregate, "  ", "  ")
		if err != nil {
			fmt.Printf("  (error formatting aggregate: %v)\n", err)
		} else {
			fmt.Printf("  %s\n", string(aggJSON))
		}
	}
}
