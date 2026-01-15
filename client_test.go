package findapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_RateLimitHandling(t *testing.T) {
	retryCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle auth endpoint
		if r.URL.Path == "/auth/v1/generate" {
			resp := struct {
				AccessToken string `json:"access_token"`
				TokenType   string `json:"token_type"`
				ExpiresIn   int    `json:"expires_in"`
				Exp         int64  `json:"exp"`
				Iat         int64  `json:"iat"`
			}{
				AccessToken: "test-token",
				TokenType:   "Bearer",
				ExpiresIn:   600,
				Exp:         time.Now().Add(10 * time.Minute).Unix(),
				Iat:         time.Now().Unix(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		retryCount++
		if retryCount < 2 {
			// First request returns rate limit error
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		// Second request succeeds
		resp := struct {
			Blocks []struct {
				Height uint64 `json:"height"`
				ID     string `json:"id"`
			} `json:"blocks"`
		}{
			Blocks: []struct {
				Height uint64 `json:"height"`
				ID     string `json:"id"`
			}{
				{Height: 96708412, ID: "abc123"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient("test", "test", WithBaseURL(server.URL))

	ctx := context.Background()
	_, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
	if err != nil {
		t.Fatalf("GetBlocks failed after retry: %v", err)
	}

	if retryCount != 2 {
		t.Errorf("Expected 2 requests (1 retry), got %d", retryCount)
	}
}

func TestClient_RateLimitExhausted(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle auth endpoint
		if r.URL.Path == "/auth/v1/generate" {
			resp := struct {
				AccessToken string `json:"access_token"`
				TokenType   string `json:"token_type"`
				ExpiresIn   int    `json:"expires_in"`
				Exp         int64  `json:"exp"`
				Iat         int64  `json:"iat"`
			}{
				AccessToken: "test-token",
				TokenType:   "Bearer",
				ExpiresIn:   600,
				Exp:         time.Now().Add(10 * time.Minute).Unix(),
				Iat:         time.Now().Unix(),
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		// Always return rate limit error
		w.Header().Set("Retry-After", "1")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := NewClient("test", "test", WithBaseURL(server.URL))

	ctx := context.Background()
	_, err := client.Simple.GetBlocks().Height(96708412).Do(ctx)
	if err == nil {
		t.Fatal("Expected error after exhausting retries")
	}

	if !IsRateLimitError(err) {
		t.Errorf("Expected RateLimitError, got %T: %v", err, err)
	}
}
