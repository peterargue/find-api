package simple

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// mockClient implements the Client interface for testing
type mockClient struct {
	server *httptest.Server
}

func (m *mockClient) DoRequest(ctx context.Context, method, path string, query url.Values) (*http.Response, error) {
	u, err := url.Parse(m.server.URL + path)
	if err != nil {
		return nil, err
	}
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, method, u.String(), nil)
	if err != nil {
		return nil, err
	}

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

func TestSimpleService_GetBlocks(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method and path
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/simple/v1/blocks" {
			t.Errorf("Expected path /simple/v1/blocks, got %s", r.URL.Path)
		}

		// Check query parameters
		height := r.URL.Query().Get("height")
		if height != "96708412" {
			t.Errorf("Expected height 96708412, got %s", height)
		}

		// Mock response
		resp := BlocksResponse{
			Blocks: []Block{
				{
					Height:    96708412,
					ID:        "abc123",
					Timestamp: "2024-01-15T10:00:00Z",
					TxCount:   5,
					Transactions: []TransactionID{
						{ID: "tx1"},
						{ID: "tx2"},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create service with mock client
	client := &mockClient{server: server}
	service := NewService(client)

	// Test GetBlocks
	ctx := context.Background()
	result, err := service.GetBlocks().Height(96708412).Do(ctx)
	if err != nil {
		t.Fatalf("GetBlocks failed: %v", err)
	}

	if len(result.Blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result.Blocks))
	}

	block := result.Blocks[0]
	if block.Height != 96708412 {
		t.Errorf("Expected height 96708412, got %d", block.Height)
	}
	if block.TxCount != 5 {
		t.Errorf("Expected 5 transactions, got %d", block.TxCount)
	}
	if len(block.Transactions) != 2 {
		t.Errorf("Expected 2 transaction IDs, got %d", len(block.Transactions))
	}
}

func TestSimpleService_GetEvents(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/simple/v1/events" {
			t.Errorf("Expected path /simple/v1/events, got %s", r.URL.Path)
		}

		// Check query parameters
		name := r.URL.Query().Get("name")
		fromHeight := r.URL.Query().Get("from_height")
		toHeight := r.URL.Query().Get("to_height")

		if name != "A.test.Event" {
			t.Errorf("Expected name A.test.Event, got %s", name)
		}
		if fromHeight != "100" {
			t.Errorf("Expected from_height 100, got %s", fromHeight)
		}
		if toHeight != "200" {
			t.Errorf("Expected to_height 200, got %s", toHeight)
		}

		// Mock response
		resp := EventsResponse{
			Events: []Event{
				{
					BlockHeight:     150,
					EventIndex:      0,
					Name:            "A.test.Event",
					Timestamp:       "2024-01-15T10:00:00Z",
					TransactionHash: "tx123",
					Fields: map[string]interface{}{
						"amount": 100.0,
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetEvents().Name("A.test.Event").FromHeight(100).ToHeight(200).Do(ctx)
	if err != nil {
		t.Fatalf("GetEvents failed: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}

	event := result.Events[0]
	if event.BlockHeight != 150 {
		t.Errorf("Expected block height 150, got %d", event.BlockHeight)
	}
	if event.Name != "A.test.Event" {
		t.Errorf("Expected name A.test.Event, got %s", event.Name)
	}
}

func TestSimpleService_GetTransaction(t *testing.T) {
	txID := "b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/simple/v1/transaction" {
			t.Errorf("Expected path /simple/v1/transaction, got %s", r.URL.Path)
		}

		id := r.URL.Query().Get("id")
		if id != txID {
			t.Errorf("Expected id %s, got %s", txID, id)
		}

		resp := TransactionsResponse{
			Transactions: []Transaction{
				{
					ID:          txID,
					BlockHeight: 96708412,
					Status:      "sealed",
					Payer:       "0x1234",
					Proposer:    "0x1234",
					GasLimit:    1000,
					GasUsed:     500,
					Fee:         0.00001,
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetTransaction().ID(txID).Do(ctx)
	if err != nil {
		t.Fatalf("GetTransaction failed: %v", err)
	}

	if len(result.Transactions) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Transactions))
	}

	tx := result.Transactions[0]
	if tx.ID != txID {
		t.Errorf("Expected ID %s, got %s", txID, tx.ID)
	}
	if tx.Status != "sealed" {
		t.Errorf("Expected status sealed, got %s", tx.Status)
	}
}

func TestSimpleService_GetTransactionEvents(t *testing.T) {
	txID := "b03b47104a675dd2d594a8dd85cdc313586678f508fe67c4de0604f0a4920562"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/simple/v1/transaction/events" {
			t.Errorf("Expected path /simple/v1/transaction/events, got %s", r.URL.Path)
		}

		transactionID := r.URL.Query().Get("transaction_id")
		if transactionID != txID {
			t.Errorf("Expected transaction_id %s, got %s", txID, transactionID)
		}

		resp := TransactionEventsResponse{
			Events: []SimpleEvent{
				{
					EventIndex: 0,
					Name:       "flow.AccountCreated",
					Fields: map[string]interface{}{
						"address": "0x1234",
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetTransactionEvents().TransactionID(txID).Do(ctx)
	if err != nil {
		t.Fatalf("GetTransactionEvents failed: %v", err)
	}

	if len(result.Events) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Events))
	}

	event := result.Events[0]
	if event.Name != "flow.AccountCreated" {
		t.Errorf("Expected name flow.AccountCreated, got %s", event.Name)
	}
}

func TestService_WithOffset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/simple/v1/blocks" {
			t.Errorf("Expected path /simple/v1/blocks, got %s", r.URL.Path)
		}

		// Check offset parameter
		offset := r.URL.Query().Get("offset")
		if offset != "10" {
			t.Errorf("Expected offset 10, got %s", offset)
		}

		resp := BlocksResponse{
			Blocks: []Block{
				{
					Height:  96708412,
					ID:      "abc123",
					TxCount: 5,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetBlocks().Height(96708412).Offset(10).Do(ctx)
	if err != nil {
		t.Fatalf("GetBlocks failed: %v", err)
	}

	if len(result.Blocks) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result.Blocks))
	}
}
