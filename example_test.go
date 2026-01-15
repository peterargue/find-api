package findapi_test

import (
	"context"
	"fmt"
	"log"

	findapi "github.com/peterargue/find-api"
)

func Example_basicUsage() {
	// Create a new client with username and password
	client := findapi.NewClient("username", "password")

	// The client automatically handles JWT token generation and refresh
	ctx := context.Background()

	// Example 1: Get blocks at a specific height
	blocksResp, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}
	for _, block := range blocksResp.Blocks {
		fmt.Printf("Block %d has %d transactions\n", block.Height, block.TxCount)
	}

	// Example 2: Get events within a block range
	eventsResp, err := client.Simple.GetEvents().
		Name("A.921ea449dffec68a.FlovatarMarketplace.FlovatarPriceChanged").
		FromHeight(102968960).
		ToHeight(103850311).
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}
	for _, event := range eventsResp.Events {
		fmt.Printf("Event at block %d: %s\n", event.BlockHeight, event.Name)
	}

	// Example 3: Get a specific transaction
	txResp, err := client.Simple.GetTransaction().
		ID("b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562").
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get transaction: %v", err)
	}
	for _, tx := range txResp.Transactions {
		fmt.Printf("Transaction %s status: %s\n", tx.ID, tx.Status)
	}

	// Example 4: Get events for a specific transaction
	txEventsResp, err := client.Simple.GetTransactionEvents().
		TransactionID("b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562").
		Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get transaction events: %v", err)
	}
	for _, event := range txEventsResp.Events {
		fmt.Printf("Event: %s\n", event.Name)
	}
}

func Example_customConfiguration() {
	// Create a client with custom options
	client := findapi.NewClient(
		"username",
		"password",
		findapi.WithBaseURL("https://custom-api.example.com"),
	)

	ctx := context.Background()
	blocksResp, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
	if err != nil {
		log.Fatalf("Failed to get blocks: %v", err)
	}
	fmt.Printf("Got %d blocks\n", len(blocksResp.Blocks))
}

func Example_rateLimitHandling() {
	client := findapi.NewClient("username", "password")
	ctx := context.Background()

	// The client automatically handles rate limiting with retries
	blocksResp, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
	if err != nil {
		// Check if it's a rate limit error after exhausting retries
		if findapi.IsRateLimitError(err) {
			log.Printf("Rate limited: %v", err)
			// Handle rate limit error (e.g., implement exponential backoff)
		} else {
			log.Fatalf("Failed to get blocks: %v", err)
		}
		return
	}

	fmt.Printf("Got %d blocks\n", len(blocksResp.Blocks))
}

func Example_pagination() {
	client := findapi.NewClient("username", "password")
	ctx := context.Background()

	// Paginate through events
	offset := 0
	for {
		eventsResp, err := client.Simple.GetEvents().
			Name("A.921ea449dffec68a.FlovatarMarketplace.FlovatarPriceChanged").
			FromHeight(102968960).
			ToHeight(103850311).
			Offset(offset).
			Do(ctx)
		if err != nil {
			log.Fatalf("Failed to get events: %v", err)
		}

		if len(eventsResp.Events) == 0 {
			break
		}

		for _, event := range eventsResp.Events {
			fmt.Printf("Event at block %d: %s\n", event.BlockHeight, event.Name)
		}

		// Increment offset for next page
		offset += len(eventsResp.Events)

		// Check if we got less than the max (100), indicating last page
		if len(eventsResp.Events) < 100 {
			break
		}
	}
}
