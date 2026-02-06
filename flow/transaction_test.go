package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/transaction" {
			t.Errorf("Expected path /flow/v1/transaction, got %s", r.URL.Path)
		}

		resp := TransactionsResponse{
			Data: []Transaction{
				{
					ID:          "abc123",
					BlockHeight: 96708412,
					Payer:       "0x1234",
					Status:      "SEALED",
					GasUsed:     100,
					Fee:         0.001,
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
	result, err := service.GetTransactions().Do(ctx)
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.ID != "abc123" {
		t.Errorf("Expected ID abc123, got %s", tx.ID)
	}
	if tx.BlockHeight != 96708412 {
		t.Errorf("Expected block height 96708412, got %d", tx.BlockHeight)
	}
}

func TestFlowService_GetTransactionsWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		payer := r.URL.Query().Get("payer")
		status := r.URL.Query().Get("status")
		height := r.URL.Query().Get("height")
		includeEvents := r.URL.Query().Get("include_events")

		if payer != "0x1234" {
			t.Errorf("Expected payer 0x1234, got %s", payer)
		}
		if status != "SEALED" {
			t.Errorf("Expected status SEALED, got %s", status)
		}
		if height != "96708412" {
			t.Errorf("Expected height 96708412, got %s", height)
		}
		if includeEvents != "true" {
			t.Errorf("Expected include_events true, got %s", includeEvents)
		}

		resp := TransactionsResponse{
			Data: []Transaction{
				{ID: "abc123", Payer: payer, Status: status},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetTransactions().
		Payer("0x1234").
		Status("SEALED").
		Height(96708412).
		IncludeEvents(true).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}
}

func TestFlowService_GetTransaction(t *testing.T) {
	txID := "abc123def456"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/transaction/%s", txID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := TransactionResponse{
			Data: []TransactionDetails{
				{
					ID:          txID,
					BlockHeight: 96708412,
					Payer:       "0x1234",
					Status:      "SEALED",
					GasUsed:     100,
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

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.ID != txID {
		t.Errorf("Expected ID %s, got %s", txID, tx.ID)
	}
}

func TestFlowService_GetScheduledTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/scheduled-transaction" {
			t.Errorf("Expected path /flow/v1/scheduled-transaction, got %s", r.URL.Path)
		}

		completed := r.URL.Query().Get("completed")
		owner := r.URL.Query().Get("owner")

		if completed != "true" {
			t.Errorf("Expected completed true, got %s", completed)
		}
		if owner != "0x1234" {
			t.Errorf("Expected owner 0x1234, got %s", owner)
		}

		resp := ScheduledTransactionsResponse{
			Data: []ScheduledTransaction{
				{
					ID:          "sched123",
					Owner:       "0x1234",
					IsCompleted: true,
					Status:      "completed",
					Handler:     "TestHandler",
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
	result, err := service.GetScheduledTransactions().
		Completed(true).
		Owner("0x1234").
		Do(ctx)
	if err != nil {
		t.Fatalf("GetScheduledTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 scheduled transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.ID != "sched123" {
		t.Errorf("Expected ID sched123, got %s", tx.ID)
	}
	if !tx.IsCompleted {
		t.Error("Expected transaction to be completed")
	}
}

func TestFlowService_TransactionRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetTransaction without ID
	_, err := service.GetTransaction().Do(ctx)
	if err == nil {
		t.Error("Expected error when transaction ID is not provided")
	}
}
