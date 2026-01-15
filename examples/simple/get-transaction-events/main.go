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

	// Get events for a specific transaction
	transactionID := "b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562"

	fmt.Printf("Fetching events for transaction: %s\n\n", transactionID)

	eventsResp, err := client.Simple.GetTransactionEvents().TransactionID(transactionID).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get transaction events: %v", err)
	}

	// Display event information
	fmt.Printf("Found %d event(s)\n\n", len(eventsResp.Events))

	if len(eventsResp.Events) == 0 {
		fmt.Println("No events found for this transaction")
		return
	}

	for i, event := range eventsResp.Events {
		fmt.Printf("Event #%d:\n", i+1)
		fmt.Printf("  Name: %s\n", event.Name)
		fmt.Printf("  Event Index: %d\n", event.EventIndex)

		// Pretty print the event fields
		if len(event.Fields) > 0 {
			fmt.Printf("  Fields:\n")
			fieldsJSON, err := json.MarshalIndent(event.Fields, "    ", "  ")
			if err != nil {
				fmt.Printf("    (error formatting fields: %v)\n", err)
			} else {
				fmt.Printf("    %s\n", string(fieldsJSON))
			}
		}
		fmt.Println()
	}

	// Example: Get events with pagination
	fmt.Println("--- Pagination Example ---")
	offset := 0
	totalEvents := 0

	for {
		pageEvents, err := client.Simple.GetTransactionEvents().TransactionID(transactionID).Offset(offset).Do(ctx)
		if err != nil {
			log.Fatalf("Failed to get transaction events (offset %d): %v", offset, err)
		}

		if len(pageEvents.Events) == 0 {
			break
		}

		totalEvents += len(pageEvents.Events)
		fmt.Printf("Fetched %d events at offset %d\n", len(pageEvents.Events), offset)

		offset += len(pageEvents.Events)

		// For most transactions, all events fit in one page
		// This is just to demonstrate the pagination capability
		if offset > 100 {
			break
		}
	}

	fmt.Printf("\nTotal events fetched: %d\n", totalEvents)

	// Example: Filter events by name
	fmt.Println("\n--- Filtering Events by Name ---")
	eventsResp, err = client.Simple.GetTransactionEvents().TransactionID(transactionID).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get transaction events: %v", err)
	}

	// Count events by type
	eventCounts := make(map[string]int)
	for _, event := range eventsResp.Events {
		eventCounts[event.Name]++
	}

	fmt.Println("Event types in this transaction:")
	for name, count := range eventCounts {
		fmt.Printf("  %s: %d\n", name, count)
	}
}
