package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

// mockClient implements the Client interface for testing
type mockClient struct {
	server *httptest.Server
}

func (m *mockClient) DoRequestWithBasicAuth(ctx context.Context, method, path string, query url.Values, username, password string) (*http.Response, error) {
	u, err := url.Parse(m.server.URL + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

	// Add basic auth credentials as headers for test verification
	req.SetBasicAuth(username, password)

	return http.DefaultClient.Do(req)
}

func (m *mockClient) DecodeResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	if v == nil {
		return nil
	}

	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

func TestAuthService_GenerateToken(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/auth/v1/generate" {
			t.Errorf("Expected path /auth/v1/generate, got %s", r.URL.Path)
		}

		// Check Basic Auth
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected Basic Auth header")
		}
		if username != "testuser" {
			t.Errorf("Expected username 'testuser', got %s", username)
		}
		if password != "testpass" {
			t.Errorf("Expected password 'testpass', got %s", password)
		}

		// Check query parameters
		expiry := r.URL.Query().Get("expiry")
		if expiry != "10m0s" {
			t.Errorf("Expected expiry '10m0s', got %s", expiry)
		}

		// Mock response
		resp := TokenResponse{
			AccessToken: "test-jwt-token",
			TokenType:   "Bearer",
			ExpiresIn:   600,
			Exp:         time.Now().Add(10 * time.Minute).Unix(),
			Iat:         time.Now().Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create service with mock client
	client := &mockClient{server: server}
	service := NewService(client, "testuser", "testpass")

	// Test GenerateToken
	ctx := context.Background()
	result, err := service.GenerateToken(ctx, 10*time.Minute)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if result.AccessToken != "test-jwt-token" {
		t.Errorf("Expected token 'test-jwt-token', got %s", result.AccessToken)
	}
	if result.TokenType != "Bearer" {
		t.Errorf("Expected token type 'Bearer', got %s", result.TokenType)
	}
	if result.ExpiresIn != 600 {
		t.Errorf("Expected expires_in 600, got %d", result.ExpiresIn)
	}
}

func TestAuthService_GenerateToken_CustomExpiry(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that expiry is set to 1 hour
		expiry := r.URL.Query().Get("expiry")
		if expiry != "1h0m0s" {
			t.Errorf("Expected expiry '1h0m0s', got %s", expiry)
		}

		// Mock response
		resp := TokenResponse{
			AccessToken: "test-jwt-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			Exp:         time.Now().Add(1 * time.Hour).Unix(),
			Iat:         time.Now().Unix(),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client, "testuser", "testpass")

	ctx := context.Background()
	result, err := service.GenerateToken(ctx, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if result.AccessToken != "test-jwt-token" {
		t.Errorf("Expected token 'test-jwt-token', got %s", result.AccessToken)
	}
	if result.ExpiresIn != 3600 {
		t.Errorf("Expected expires_in 3600, got %d", result.ExpiresIn)
	}
}

func TestAuthService_GenerateToken_Error(t *testing.T) {
	// Create mock server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Invalid credentials"))
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client, "wronguser", "wrongpass")

	ctx := context.Background()
	_, err := service.GenerateToken(ctx, 10*time.Minute)
	if err == nil {
		t.Fatal("Expected error for invalid credentials")
	}
}
