package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	findapi "github.com/peterargue/find-api"
)

func main() {
	// Load .env if it exists; ignore error if file is absent.
	_ = godotenv.Load()

	apiURL := os.Getenv("FIND_API_URL")
	if apiURL == "" {
		apiURL = findapi.FindApiURL
	}

	var client *findapi.Client

	user := os.Getenv("FIND_API_USER")
	pass := os.Getenv("FIND_API_PASS")

	if user != "" && pass != "" {
		client = findapi.NewClient(user, pass, findapi.WithBaseURL(apiURL))
	} else {
		// Fall back to the stored CLI token.
		token, exp, err := loadStoredToken()
		if err != nil {
			log.Fatalf("No credentials available: set FIND_API_USER/FIND_API_PASS or run 'find auth login'. (%v)", err)
		}
		client = findapi.NewClient("", "", findapi.WithBaseURL(apiURL), findapi.WithToken(token, exp))
	}

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

type tokenFile struct {
	AccessToken string `json:"access_token"`
	Exp         int64  `json:"exp"`
}

func loadStoredToken() (string, int64, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".config", "find-cli", "token.json")
	b, err := os.ReadFile(path)
	if err != nil {
		return "", 0, err
	}
	var tf tokenFile
	if err := json.Unmarshal(b, &tf); err != nil {
		return "", 0, err
	}
	if time.Now().Add(time.Minute).After(time.Unix(tf.Exp, 0)) {
		return "", 0, os.ErrDeadlineExceeded
	}
	return tf.AccessToken, tf.Exp, nil
}
