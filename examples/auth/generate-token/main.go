package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

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

	// Generate a JWT token with 1 hour expiry
	fmt.Println("Generating JWT token with 1 hour expiry...")
	token, err := client.Auth.GenerateToken(ctx, time.Hour)
	if err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	// Display token information
	fmt.Printf("\nToken generated successfully!\n")
	fmt.Printf("Access Token: %s\n", token.AccessToken)
	fmt.Printf("Token Type: %s\n", token.TokenType)
	fmt.Printf("Expires In: %d seconds\n", token.ExpiresIn)
	fmt.Printf("Expiry Time (Unix): %d\n", token.Exp)
	fmt.Printf("Issued At (Unix): %d\n", token.Iat)

	// Convert Unix timestamps to human-readable format
	expiryTime := time.Unix(token.Exp, 0)
	issuedTime := time.Unix(token.Iat, 0)
	fmt.Printf("\nExpires At: %s\n", expiryTime.Format(time.RFC3339))
	fmt.Printf("Issued At: %s\n", issuedTime.Format(time.RFC3339))

	// Generate a token with different expiry (24 hours)
	fmt.Println("\nGenerating JWT token with 24 hour expiry...")
	longToken, err := client.Auth.GenerateToken(ctx, 24*time.Hour)
	if err != nil {
		log.Fatalf("Failed to generate long-lived token: %v", err)
	}

	longExpiryTime := time.Unix(longToken.Exp, 0)
	fmt.Printf("Long-lived token expires at: %s\n", longExpiryTime.Format(time.RFC3339))
}
