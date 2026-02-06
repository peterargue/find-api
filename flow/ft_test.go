package flow

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

func TestFlowService_GetFTs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/ft" {
			t.Errorf("Expected path /flow/v1/ft, got %s", r.URL.Path)
		}

		resp := FTListResponse{
			Data: []FungibleToken{
				{
					Address:      "0x1654653399040a61",
					ContractName: "FlowToken",
					Symbol:       "FLOW",
					Name:         "Flow Token",
					Decimals:     8,
					Token:        "A.1654653399040a61.FlowToken.Vault",
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
	result, err := service.GetFTs().Do(ctx)
	if err != nil {
		t.Fatalf("GetFTs failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}

	token := result.Data[0]
	if token.Symbol != "FLOW" {
		t.Errorf("Expected symbol FLOW, got %s", token.Symbol)
	}
}

func TestFlowService_GetFT(t *testing.T) {
	tokenID := "A.1654653399040a61.FlowToken.Vault"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/ft/%s", tokenID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := FungibleTokenResponse{
			Data: []FungibleTokenDetails{
				{
					FungibleToken: FungibleToken{
						Address:      "0x1654653399040a61",
						ContractName: "FlowToken",
						Symbol:       "FLOW",
						Name:         "Flow Token",
						Decimals:     8,
						Token:        tokenID,
					},
					Stats: FTStats{
						OwnerCounts:  1000,
						TotalBalance: 1000000.0,
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
	result, err := service.GetFT().Token(tokenID).Do(ctx)
	if err != nil {
		t.Fatalf("GetFT failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}

	token := result.Data[0]
	if token.Symbol != "FLOW" {
		t.Errorf("Expected symbol FLOW, got %s", token.Symbol)
	}
	if token.Stats.OwnerCounts != 1000 {
		t.Errorf("Expected 1000 owners, got %d", token.Stats.OwnerCounts)
	}
}

func TestFlowService_GetFTTransfers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/ft/transfer" {
			t.Errorf("Expected path /flow/v1/ft/transfer, got %s", r.URL.Path)
		}

		token := r.URL.Query().Get("token")
		if token != "A.1654653399040a61.FlowToken.Vault" {
			t.Errorf("Expected token A.1654653399040a61.FlowToken.Vault, got %s", token)
		}

		resp := TransfersResponse{
			Data: []FTTransfer{
				{
					Amount:        100.5,
					BlockHeight:   96708412,
					Direction:     "withdraw",
					Sender:        "0x1234",
					Receiver:      "0x5678",
					Token:         "A.1654653399040a61.FlowToken.Vault",
					TransactionID: "abc123",
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
	result, err := service.GetFTTransfers().
		Token("A.1654653399040a61.FlowToken.Vault").
		Do(ctx)
	if err != nil {
		t.Fatalf("GetFTTransfers failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transfer, got %d", len(result.Data))
	}

	transfer := result.Data[0]
	if transfer.Amount != 100.5 {
		t.Errorf("Expected amount 100.5, got %f", transfer.Amount)
	}
	if transfer.Direction != "withdraw" {
		t.Errorf("Expected direction withdraw, got %s", transfer.Direction)
	}
}

func TestFlowService_GetFTHoldings(t *testing.T) {
	tokenID := "A.1654653399040a61.FlowToken.Vault"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/ft/%s/holding", tokenID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := FTHoldingResponse{
			Data: []FTHolding{
				{
					Address:    "0x1234",
					Balance:    1000.0,
					Percentage: 10.5,
					Token:      tokenID,
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
	result, err := service.GetFTHoldings().Token(tokenID).Do(ctx)
	if err != nil {
		t.Fatalf("GetFTHoldings failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 holding, got %d", len(result.Data))
	}

	holding := result.Data[0]
	if holding.Balance != 1000.0 {
		t.Errorf("Expected balance 1000.0, got %f", holding.Balance)
	}
	if holding.Percentage != 10.5 {
		t.Errorf("Expected percentage 10.5, got %f", holding.Percentage)
	}
}

func TestFlowService_GetFTAccountToken(t *testing.T) {
	tokenID := "A.1654653399040a61.FlowToken.Vault"
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/ft/%s/account/%s", tokenID, address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := AccountFungibleTokenResponse{
			Data: []Vault{
				{
					Address:     address,
					Balance:     500.0,
					BlockHeight: 96708412,
					Token:       tokenID,
					VaultID:     1,
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
	result, err := service.GetFTAccountToken().
		Token(tokenID).
		Address(address).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetFTAccountToken failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 vault, got %d", len(result.Data))
	}

	vault := result.Data[0]
	if vault.Balance != 500.0 {
		t.Errorf("Expected balance 500.0, got %f", vault.Balance)
	}
	if vault.Address != address {
		t.Errorf("Expected address %s, got %s", address, vault.Address)
	}
}

func TestFlowService_GetFTsWithPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		if limit != "10" {
			t.Errorf("Expected limit 10, got %s", limit)
		}
		if offset != "5" {
			t.Errorf("Expected offset 5, got %s", offset)
		}

		resp := FTListResponse{
			Data: []FungibleToken{
				{Symbol: "FLOW"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetFTs().Limit(10).Offset(5).Do(ctx)
	if err != nil {
		t.Fatalf("GetFTs failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 token, got %d", len(result.Data))
	}
}

func TestFlowService_RequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetFT without token
	_, err := service.GetFT().Do(ctx)
	if err == nil {
		t.Error("Expected error when token is not provided")
	}

	// Test GetFTHoldings without token
	_, err = service.GetFTHoldings().Do(ctx)
	if err == nil {
		t.Error("Expected error when token is not provided")
	}

	// Test GetFTAccountToken without token
	_, err = service.GetFTAccountToken().Address("0x1234").Do(ctx)
	if err == nil {
		t.Error("Expected error when token is not provided")
	}

	// Test GetFTAccountToken without address
	_, err = service.GetFTAccountToken().Token("token").Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}
}
