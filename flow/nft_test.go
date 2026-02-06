package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetNFTCollections(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/nft" {
			t.Errorf("Expected path /flow/v1/nft, got %s", r.URL.Path)
		}

		resp := NFTCollectionResponse{
			Data: []NFTCollection{
				{
					NFTType: "A.0b2a3299cc857e29.TopShot.NFT",
					Name:    "NBA Top Shot",
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
	result, err := service.GetNFTCollections().Do(ctx)
	if err != nil {
		t.Fatalf("GetNFTCollections failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(result.Data))
	}

	collection := result.Data[0]
	if collection.Name != "NBA Top Shot" {
		t.Errorf("Expected name 'NBA Top Shot', got %s", collection.Name)
	}
}

func TestFlowService_GetNFTCollection(t *testing.T) {
	nftType := "A.0b2a3299cc857e29.TopShot.NFT"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/nft/%s", nftType)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := NFTCollectionDetailsResponse{
			Data: []NFTCollectionDetails{
				{
					NFTCollection: NFTCollection{
						NFTType: nftType,
						Name:    "NBA Top Shot",
					},
					ItemCount: 10000,
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
	result, err := service.GetNFTCollection().NFTType(nftType).Do(ctx)
	if err != nil {
		t.Fatalf("GetNFTCollection failed: %v", err)
	}

	if len(result.Data) == 0 {
		t.Fatal("Expected at least 1 collection")
	}

	if result.Data[0].NFTType != nftType {
		t.Errorf("Expected nft_type %s, got %s", nftType, result.Data[0].NFTType)
	}
}

func TestFlowService_GetNFTTransfers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/nft/transfer" {
			t.Errorf("Expected path /flow/v1/nft/transfer, got %s", r.URL.Path)
		}

		height := r.URL.Query().Get("height")
		if height != "96708412" {
			t.Errorf("Expected height 96708412, got %s", height)
		}

		resp := NFTTransfersResponse{
			Data: []NFTTransfer{
				{
					BlockHeight:   96708412,
					Sender:        "0x1234",
					Receiver:      "0x5678",
					NFTType:       "A.0b2a3299cc857e29.TopShot.NFT",
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
	result, err := service.GetNFTTransfers().Height(96708412).Do(ctx)
	if err != nil {
		t.Fatalf("GetNFTTransfers failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 transfer, got %d", len(result.Data))
	}

	transfer := result.Data[0]
	if transfer.BlockHeight != 96708412 {
		t.Errorf("Expected block height 96708412, got %d", transfer.BlockHeight)
	}
}

func TestFlowService_GetAccountNFTCollections(t *testing.T) {
	address := "0x1654653399040a61"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/nft", address)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := AccountNFTCollectionsResponse{
			Data: []AccountNFTCollection{
				{
					Owner:    address,
					NFTType:  "A.0b2a3299cc857e29.TopShot.NFT",
					Name:     "NBA Top Shot",
					NFTCount: 5,
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
	result, err := service.GetAccountNFTCollections().Address(address).Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountNFTCollections failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 collection, got %d", len(result.Data))
	}

	collection := result.Data[0]
	if collection.NFTCount != 5 {
		t.Errorf("Expected NFT count 5, got %d", collection.NFTCount)
	}
}

func TestFlowService_GetAccountNFTs(t *testing.T) {
	address := "0x1654653399040a61"
	nftType := "A.0b2a3299cc857e29.TopShot.NFT"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/account/%s/nft/%s", address, nftType)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		validOnly := r.URL.Query().Get("valid_only")
		if validOnly != "true" {
			t.Errorf("Expected valid_only=true, got %s", validOnly)
		}

		resp := AccountNFTResponse{
			Data: []AccountNFT{
				{
					ID:      "123",
					Name:    "Moment #123",
					NFTType: "A.0b2a3299cc857e29.TopShot.NFT",
					Address: address,
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
	result, err := service.GetAccountNFTs().
		Address(address).
		NFTType(nftType).
		ValidOnly(true).
		Do(ctx)
	if err != nil {
		t.Fatalf("GetAccountNFTs failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 NFT, got %d", len(result.Data))
	}

	nft := result.Data[0]
	if nft.Name != "Moment #123" {
		t.Errorf("Expected name 'Moment #123', got %s", nft.Name)
	}
}

func TestFlowService_NFTRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetNFTCollection without nft_type
	_, err := service.GetNFTCollection().Do(ctx)
	if err == nil {
		t.Error("Expected error when nft_type is not provided")
	}

	// Test GetAccountNFTCollections without address
	_, err = service.GetAccountNFTCollections().Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetAccountNFTs without address
	_, err = service.GetAccountNFTs().NFTType("A.0b2a3299cc857e29.TopShot.NFT").Do(ctx)
	if err == nil {
		t.Error("Expected error when address is not provided")
	}

	// Test GetAccountNFTs without nft_type
	_, err = service.GetAccountNFTs().Address("0x1234").Do(ctx)
	if err == nil {
		t.Error("Expected error when nft_type is not provided")
	}
}
