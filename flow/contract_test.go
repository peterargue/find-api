package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetContracts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/contract" {
			t.Errorf("Expected path /flow/v1/contract, got %s", r.URL.Path)
		}

		resp := ContractResponse{
			Data: []Contract{
				{
					Address:      "0x1654653399040a61",
					BlockHeight:  7601063,
					ContractName: "FlowToken",
					Identifier:   "A.1654653399040a61.FlowToken",
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
	result, err := service.GetContracts().Do(ctx)
	if err != nil {
		t.Fatalf("GetContracts failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(result.Data))
	}

	contract := result.Data[0]
	if contract.ContractName != "FlowToken" {
		t.Errorf("Expected contract name FlowToken, got %s", contract.ContractName)
	}
	if contract.BlockHeight != 7601063 {
		t.Errorf("Expected block height 7601063, got %d", contract.BlockHeight)
	}
}

func TestFlowService_GetContractsByIdentifier(t *testing.T) {
	identifier := "A.1654653399040a61.FlowToken"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/contract/%s", identifier)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := ContractResponse{
			Data: []Contract{
				{
					Address:      "0x1654653399040a61",
					BlockHeight:  7601063,
					ContractName: "FlowToken",
					Identifier:   identifier,
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
	result, err := service.GetContractsByIdentifier().Identifier(identifier).Do(ctx)
	if err != nil {
		t.Fatalf("GetContractsByIdentifier failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(result.Data))
	}

	contract := result.Data[0]
	if contract.Identifier != identifier {
		t.Errorf("Expected identifier %s, got %s", identifier, contract.Identifier)
	}
}

func TestFlowService_GetContract(t *testing.T) {
	identifier := "A.1654653399040a61.FlowToken"
	id := "0x1654653399040a61"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/contract/%s/%s", identifier, id)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := ContractResponse{
			Data: []Contract{
				{
					Address:      id,
					BlockHeight:  7601063,
					ContractName: "FlowToken",
					Identifier:   identifier,
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
	result, err := service.GetContract().Identifier(identifier).ID(id).Do(ctx)
	if err != nil {
		t.Fatalf("GetContract failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(result.Data))
	}

	contract := result.Data[0]
	if contract.Address != id {
		t.Errorf("Expected address %s, got %s", id, contract.Address)
	}
	if contract.Identifier != identifier {
		t.Errorf("Expected identifier %s, got %s", identifier, contract.Identifier)
	}
}

func TestFlowService_GetContractsWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")

		if limit != "50" {
			t.Errorf("Expected limit 50, got %s", limit)
		}
		if offset != "10" {
			t.Errorf("Expected offset 10, got %s", offset)
		}

		resp := ContractResponse{
			Data: []Contract{
				{Address: "0x1234", ContractName: "TestContract"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetContracts().
		Limit(50).
		Offset(10).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetContracts failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 contract, got %d", len(result.Data))
	}
}

func TestFlowService_ContractRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetContractsByIdentifier without identifier
	_, err := service.GetContractsByIdentifier().Do(ctx)
	if err == nil {
		t.Error("Expected error when identifier is not provided")
	}

	// Test GetContract without identifier
	_, err = service.GetContract().ID("0x1234").Do(ctx)
	if err == nil {
		t.Error("Expected error when identifier is not provided")
	}

	// Test GetContract without ID
	_, err = service.GetContract().Identifier("A.1234.Test").Do(ctx)
	if err == nil {
		t.Error("Expected error when ID is not provided")
	}
}
