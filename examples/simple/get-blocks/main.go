package main

import (
	"context"
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

	// Create client - JWT authentication happens automatically
	client := findapi.NewClient(username, password)

	ctx := context.Background()

	// Get blocks at a specific height
	height := 96708412
	fmt.Printf("Fetching blocks at height %d...\n\n", height)

	blocks, err := client.Simple.GetBlocks().
		Height(height).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}

	// Display block information
	fmt.Printf("Found %d block(s)\n\n", len(blocks.Blocks))

	for i, block := range blocks.Blocks {
		fmt.Printf("Block #%d:\n", i+1)
		fmt.Printf("  Height: %d\n", block.Height)
		fmt.Printf("  ID: %s\n", block.ID)
		fmt.Printf("  Timestamp: %s\n", block.Timestamp)
		fmt.Printf("  Transaction Count: %d\n", block.TxCount)
		fmt.Printf("  Transactions:\n")
		for j, tx := range block.Transactions {
			fmt.Printf("    %d. %s\n", j+1, tx.ID)
		}
		fmt.Println()
	}

	// Example: Get blocks with offset
	offset := 5
	fmt.Printf("Fetching blocks at height %d with offset %d...\n\n", height, offset)

	blocksWithOffset, err := client.Simple.GetBlocks().Height(height).Offset(offset).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get blocks with offset: %v", err)
	}

	fmt.Printf("Found %d block(s) with offset\n", len(blocksWithOffset.Blocks))
	for i, block := range blocksWithOffset.Blocks {
		fmt.Printf("  %d. Height: %d, ID: %s\n", i+1, block.Height, block.ID)
	}
}
