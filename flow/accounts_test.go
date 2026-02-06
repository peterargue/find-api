package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetAccounts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/account" {
			t.Errorf("Expected path /flow/v1/account, got %s", r.URL.Path)
		}

		resp := AccountsResponse{
			Data: []Account{
				{
					Address:     "0x1234",
					FlowBalance: 100.5,
					Height:      96708412,
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
	result, err := service.GetAccounts().Do(ctx)
	if err != nil {
		t.Fatalf("GetAccounts failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 account, got %d", len(result.Data))
	}

	account := result.Data[0]
	if account.Address != "0x1234" {
		t.Errorf("Expected address 0x1234, got %s", account.Address)
	}
	if account.FlowBalance != 100.5 {
		t.Errorf("Expected balance 100.5, got %f", account.FlowBalance)
	}
}

func TestFlowService_GetAccount(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := AccountDetailsResponse{
			Data: []CombinedAccountDetails{
				{
					Address:     address,
					FlowBalance: 100.5,
					Keys: []KeyInfo{
						{
							Index:     0,
							PublicKey: "abc123",
						},
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
	result, err := service.GetAccount().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccount failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 account, got %d", len(result.Data))
	}

	account := result.Data[0]
	if account.Address != address {
		t.Errorf("Expected address %s, got %s", address, account.Address)
	}
	if len(account.Keys) != 1 {
		t.Errorf("Expected 1 key, got %d", len(account.Keys))
	}
}

func TestFlowService_GetAccountFTs(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/ft", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := AccountFTCollectionsResponse{
			Data: []AccountFTCollection{
				{
					Address: address,
					Balance: "100.5",
					Token:   "A.1654653399040a61.FlowToken.Vault",
					VaultID: 1,
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
	result, err := service.GetAccountFTs().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountFTs failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 FT collection, got %d", len(result.Data))
	}

	ft := result.Data[0]
	if ft.Balance != "100.5" {
		t.Errorf("Expected balance 100.5, got %s", ft.Balance)
	}
}

func TestFlowService_GetAccountFTHoldings(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/ft/holding", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := FTHoldingResponse{
			Data: []FTHolding{
				{
					Address:    address,
					Balance:    100.5,
					Percentage: 10.0,
					Token:      "A.1654653399040a61.FlowToken.Vault",
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
	result, err := service.GetAccountFTHoldings().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountFTHoldings failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 holding, got %d", len(result.Data))
	}

	holding := result.Data[0]
	if holding.Percentage != 10.0 {
		t.Errorf("Expected percentage 10.0, got %f", holding.Percentage)
	}
}

func TestFlowService_GetAccountFTTransfers(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/ft/transfer", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := TransfersResponse{
			Data: []FTTransfer{
				{
					Amount:      100.5,
					BlockHeight: 96708412,
					Direction:   "withdraw",
					Sender:      address,
					Receiver:    "0x5678",
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
	result, err := service.GetAccountFTTransfers().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountFTTransfers failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transfer, got %d", len(result.Data))
	}

	transfer := result.Data[0]
	if transfer.Amount != 100.5 {
		t.Errorf("Expected amount 100.5, got %f", transfer.Amount)
	}
}

func TestFlowService_GetAccountFTToken(t *testing.T) {
	address := "0x1234"
	token := "A.1654653399040a61.FlowToken.Vault"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/ft/%s", address, token)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := AccountFungibleTokenResponse{
			Data: []Vault{
				{
					Address:     address,
					Balance:     100.5,
					BlockHeight: 96708412,
					Token:       token,
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
	result, err := service.GetAccountFTToken().Address(address).Token(token).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountFTToken failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 vault, got %d", len(result.Data))
	}

	vault := result.Data[0]
	if vault.Balance != 100.5 {
		t.Errorf("Expected balance 100.5, got %f", vault.Balance)
	}
}

func TestFlowService_GetAccountFTTokenTransfers(t *testing.T) {
	address := "0x1234"
	token := "A.1654653399040a61.FlowToken.Vault"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/ft/%s/transfer", address, token)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := TransfersResponse{
			Data: []FTTransfer{
				{
					Amount:      100.5,
					BlockHeight: 96708412,
					Direction:   "deposit",
					Token:       token,
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
	result, err := service.GetAccountFTTokenTransfers().
		Address(address).
		Token(token).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountFTTokenTransfers failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transfer, got %d", len(result.Data))
	}
}

func TestFlowService_GetAccountTaxReport(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/tax-report", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := TaxReportResponse{
			Data: []TaxReportEntry{
				{
					Address:     address,
					Amount:      100.5,
					BlockHeight: 96708412,
					Direction:   "deposit",
					Type:        "ft",
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
	result, err := service.GetAccountTaxReport().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountTaxReport failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(result.Data))
	}

	entry := result.Data[0]
	if entry.Amount != 100.5 {
		t.Errorf("Expected amount 100.5, got %f", entry.Amount)
	}
}

func TestFlowService_GetAccountTransactions(t *testing.T) {
	address := "0x1234"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/transaction", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		includeEvents := r.URL.Query().Get("include_events")
		if includeEvents != "true" {
			t.Errorf("Expected include_events=true, got %s", includeEvents)
		}

		resp := AccountTransactionsResponse{
			Data: []AccountTransaction{
				{
					TransactionID: "abc123",
					BlockHeight:   96708412,
					Status:        "sealed",
					Payer:         address,
					Fee:           0.001,
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
	result, err := service.GetAccountTransactions().
		Address(address).
		IncludeEvents(true).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountTransactions failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transaction, got %d", len(result.Data))
	}

	tx := result.Data[0]
	if tx.TransactionID != "abc123" {
		t.Errorf("Expected transaction ID abc123, got %s", tx.TransactionID)
	}
}

func TestFlowService_AccountRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetAccount without address
	_, err := service.GetAccount().Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetAccountFTs without address
	_, err = service.GetAccountFTs().Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetAccountFTToken without address
	_, err = service.GetAccountFTToken().Token("token").Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetAccountFTToken without token
	_, err = service.GetAccountFTToken().Address("0x1234").Do(ctx)
	if err == nil {
		t.Error("Expected error when token is not provided")
	}
}

// mockClient is defined in ft_test.go
