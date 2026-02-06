package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetEvmTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/evm/token" {
			t.Errorf("Expected path /flow/v1/evm/token, got %s", r.URL.Path)
		}

		resp := EvmTokenResponse{
			Data: []EvmToken{
				{
					ContractAddressHash: "0x1234567890abcdef",
					Name:                "Flow Token",
					Symbol:              "FLOW",
					Decimals:            18,
					TotalSupply:         "1000000000",
					Holders:             1000,
					Transfers:           5000,
					Type:                "ERC-20",
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
	result, err := service.GetEvmTokens().Do(ctx)
	if err != nil {
		t.Fatalf("GetEvmTokens failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}

	token := result.Data[0]
	if token.Symbol != "FLOW" {
		t.Errorf("Expected symbol FLOW, got %s", token.Symbol)
	}
	if token.Decimals != 18 {
		t.Errorf("Expected 18 decimals, got %d", token.Decimals)
	}
}

func TestFlowService_GetEvmToken(t *testing.T) {
	address := "0x1234567890abcdef"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/evm/token/%s", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := EvmTokenResponse{
			Data: []EvmToken{
				{
					ContractAddressHash: address,
					Name:                "Test Token",
					Symbol:              "TEST",
					Decimals:            18,
					TotalSupply:         "1000000",
					Type:                "ERC-20",
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
	result, err := service.GetEvmToken().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetEvmToken failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}

	token := result.Data[0]
	if token.ContractAddressHash != address {
		t.Errorf("Expected address %s, got %s", address, token.ContractAddressHash)
	}
}

func TestFlowService_GetEvmTransactions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/evm/transaction" {
			t.Errorf("Expected path /flow/v1/evm/transaction, got %s", r.URL.Path)
		}

		height := r.URL.Query().Get("height")
		if height != "96708412" {
			t.Errorf("Expected height 96708412, got %s", height)
		}

		resp := EvmTransactionResponse{
			Data: []EvmTransaction{
				{
					Hash:            "0xabc123",
					BlockNumber:     96708412,
					From:            "0x1234",
					To:              "0x5678",
					Value:           "1000000000000000000",
					GasLimit:        "21000",
					GasUsed:         "21000",
					GasPrice:        "1000000000",
					Status:          "success",
					TransactionIndex: 0,
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
	result, err := service.GetEvmTransactions().Height(96708412).Do(ctx)
	if err != nil {
		t.Fatalf("GetEvmTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.Hash != "0xabc123" {
		t.Errorf("Expected hash 0xabc123, got %s", tx.Hash)
	}
	if tx.BlockNumber != 96708412 {
		t.Errorf("Expected block number 96708412, got %d", tx.BlockNumber)
	}
}

func TestFlowService_GetEvmTransaction(t *testing.T) {
	hash := "0xabc123def456"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/evm/transaction/%s", hash)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := EvmTransaction{
			Hash:         hash,
			BlockNumber:  96708412,
			From:         "0x1234",
			To:           "0x5678",
			Value:        "1000000000000000000",
			GasLimit:     "21000",
			GasUsed:      "21000",
			Status:       "success",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetEvmTransaction().Hash(hash).Do(ctx)
	if err != nil {
		t.Fatalf("GetEvmTransaction failed: %v", err)
	}

	if result.Hash != hash {
		t.Errorf("Expected hash %s, got %s", hash, result.Hash)
	}
	if result.BlockNumber != 96708412 {
		t.Errorf("Expected block number 96708412, got %d", result.BlockNumber)
	}
}

func TestFlowService_GetEvmTokensWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		typ := r.URL.Query().Get("type")
		name := r.URL.Query().Get("name")
		limit := r.URL.Query().Get("limit")

		if typ != "ERC-20" {
			t.Errorf("Expected type ERC-20, got %s", typ)
		}
		if name != "Flow" {
			t.Errorf("Expected name Flow, got %s", name)
		}
		if limit != "50" {
			t.Errorf("Expected limit 50, got %s", limit)
		}

		resp := EvmTokenResponse{
			Data: []EvmToken{
				{Symbol: "FLOW", Type: "ERC-20"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetEvmTokens().
		Type("ERC-20").
		Name("Flow").
		Limit(50).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetEvmTokens failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}
}

func TestFlowService_EvmRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetEvmToken without address
	_, err := service.GetEvmToken().Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetEvmTransaction without hash
	_, err = service.GetEvmTransaction().Do(ctx)
	if err == nil {
		t.Error("Expected error when hash is not provided")
	}
}
