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

	// Query events by name within a block range
	eventName := "A.1654653399040a61.FlowToken.TokensWithdrawn"
	fromHeight := 85000000
	toHeight := 85000100

	fmt.Printf("Fetching events: %s\n", eventName)
	fmt.Printf("Block range: %d to %d\n\n", fromHeight, toHeight)

	events, err := client.Simple.GetEvents().
		Name(eventName).
		FromHeight(fromHeight).
		ToHeight(toHeight).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}

	// Display event information
	fmt.Printf("Found %d event(s)\n\n", len(events.Events))

	for i, event := range events.Events {
		fmt.Printf("Event #%d:\n", i+1)
		fmt.Printf("  Name: %s\n", event.Name)
		fmt.Printf("  Block Height: %d\n", event.BlockHeight)
		fmt.Printf("  Event Index: %d\n", event.EventIndex)
		fmt.Printf("  Transaction Hash: %s\n", event.TransactionHash)
		fmt.Printf("  Timestamp: %s\n", event.Timestamp)

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

		// Only show first 5 events to avoid cluttering output
		if i >= 4 {
			remaining := len(events.Events) - 5
			if remaining > 0 {
				fmt.Printf("... and %d more event(s)\n", remaining)
			}
			break
		}
	}

	// Example: Pagination through events
	fmt.Println("\n--- Pagination Example ---")
	offset := 0
	totalFetched := 0

	for {
		pageEvents, err := client.Simple.GetEvents().
			Name(eventName).
			FromHeight(fromHeight).
			ToHeight(toHeight).
			Offset(offset).
			Do(ctx)
		if err != nil {
			log.Fatalf("Failed to get events (offset %d): %v", offset, err)
		}

		if len(pageEvents.Events) == 0 {
			break
		}

		totalFetched += len(pageEvents.Events)
		fmt.Printf("Page (offset %d): fetched %d events\n", offset, len(pageEvents.Events))

		offset += len(pageEvents.Events)

		// Stop after fetching a few pages for this example
		if totalFetched >= 1000 {
			fmt.Println("Stopping pagination example after fetching 1000+ events")
			break
		}

		// Check if we got less than the max (100), indicating last page
		if len(pageEvents.Events) < 100 {
			break
		}
	}

	fmt.Printf("\nTotal events fetched across all pages: %d\n", totalFetched)
}
