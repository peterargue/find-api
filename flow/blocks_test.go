package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetBlocks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/block" {
			t.Errorf("Expected path /flow/v1/block, got %s", r.URL.Path)
		}

		resp := BlockResponse{
			Data: []Block{
				{
					Height:      96708412,
					ID:          "abc123",
					Timestamp:   "2024-01-15T10:00:00Z",
					Tx:          5,
					Fees:        0.001,
					EvmTxCount:  2,
					SurgeFactor: 1.0,
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
	result, err := service.GetBlocks().Do(ctx)
	if err != nil {
		t.Fatalf("GetBlocks failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result.Data))
	}

	block := result.Data[0]
	if block.Height != 96708412 {
		t.Errorf("Expected height 96708412, got %d", block.Height)
	}
	if block.Tx != 5 {
		t.Errorf("Expected 5 transactions, got %d", block.Tx)
	}
}

func TestFlowService_GetBlock(t *testing.T) {
	height := uint64(96708412)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/block/%d", height)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := BlockResponse{
			Data: []Block{
				{
					Height:    height,
					ID:        "abc123",
					Timestamp: "2024-01-15T10:00:00Z",
					Tx:        5,
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
	result, err := service.GetBlock().Height(height).Do(ctx)
	if err != nil {
		t.Fatalf("GetBlock failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result.Data))
	}

	block := result.Data[0]
	if block.Height != height {
		t.Errorf("Expected height %d, got %d", height, block.Height)
	}
}

func TestFlowService_GetBlockServiceEvents(t *testing.T) {
	height := uint64(96708412)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/block/%d/service-event", height)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		heightParam := r.URL.Query().Get("height")
		if heightParam != fmt.Sprintf("%d", height) {
			t.Errorf("Expected height %d, got %s", height, heightParam)
		}

		resp := BlockServiceEventResponse{
			Data: []BlockServiceEvent{
				{
					BlockHeight: height,
					EventIndex:  0,
					EventType:   "flow.EpochSetup",
					Timestamp:   "2024-01-15T10:00:00Z",
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
	result, err := service.GetBlockServiceEvents().Height(height).Do(ctx)
	if err != nil {
		t.Fatalf("GetBlockServiceEvents failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 event, got %d", len(result.Data))
	}

	event := result.Data[0]
	if event.BlockHeight != height {
		t.Errorf("Expected block height %d, got %d", height, event.BlockHeight)
	}
	if event.EventType != "flow.EpochSetup" {
		t.Errorf("Expected event type flow.EpochSetup, got %s", event.EventType)
	}
}

func TestFlowService_GetBlockTransactions(t *testing.T) {
	height := uint64(96708412)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/block/%d/transaction", height)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		includeEvents := r.URL.Query().Get("include_events")
		if includeEvents != "true" {
			t.Errorf("Expected include_events=true, got %s", includeEvents)
		}

		resp := BlockTransactionsResponse{
			Data: []BlockTransaction{
				{
					TransactionID: "abc123",
					BlockHeight:   height,
					Status:        "sealed",
					Payer:         "0x1234",
					Fee:           0.001,
					EventCount:    3,
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
	result, err := service.GetBlockTransactions().
		Height(height).
		IncludeEvents(true).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetBlockTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.BlockHeight != height {
		t.Errorf("Expected block height %d, got %d", height, tx.BlockHeight)
	}
	if tx.EventCount != 3 {
		t.Errorf("Expected 3 events, got %d", tx.EventCount)
	}
}

func TestFlowService_GetBlocksWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		height := r.URL.Query().Get("height")
		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if height != "100000000" {
			t.Errorf("Expected height 100000000, got %s", height)
		}
		if limit != "50" {
			t.Errorf("Expected limit 50, got %s", limit)
		}
		if offset != "10" {
			t.Errorf("Expected offset 10, got %s", offset)
		}

		resp := BlockResponse{
			Data: []Block{
				{Height: 99999999, ID: "block1"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetBlocks().
		Height(100000000).
		Limit(50).
		Offset(10).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetBlocks failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 block, got %d", len(result.Data))
	}
}

func TestFlowService_BlockRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetBlock without height
	_, err := service.GetBlock().Do(ctx)
	if err == nil {
		t.Error("Expected error when height is not provided")
	}

	// Test GetBlockServiceEvents without height
	_, err = service.GetBlockServiceEvents().Do(ctx)
	if err == nil {
		t.Error("Expected error when height is not provided")
	}

	// Test GetBlockTransactions without height
	_, err = service.GetBlockTransactions().Do(ctx)
	if err == nil {
		t.Error("Expected error when height is not provided")
	}
}
