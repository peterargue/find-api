package flow

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFlowService_GetNodes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/flow/v1/node" {
			t.Errorf("Expected path /flow/v1/node, got %s", r.URL.Path)
		}

		resp := NodeResponse{
			Data: []Node{
				{
					NodeID:       "abc123",
					Address:      "0x1234",
					Name:         "Test Node",
					Organization: "Test Org",
					RoleID:       1,
					Role:         "collection",
					TokensStaked: 1000000.0,
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
	result, err := service.GetNodes().Do(ctx)
	if err != nil {
		t.Fatalf("GetNodes failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 node, got %d", len(result.Data))
	}

	node := result.Data[0]
	if node.Name != "Test Node" {
		t.Errorf("Expected name 'Test Node', got %s", node.Name)
	}
	if node.RoleID != 1 {
		t.Errorf("Expected role ID 1, got %d", node.RoleID)
	}
}

func TestFlowService_GetNodesWithFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		height := r.URL.Query().Get("height")
		organization := r.URL.Query().Get("organiztion") // Note: API has typo
		roleID := r.URL.Query().Get("role_id")
		sortBy := r.URL.Query().Get("sort_by")

		if height != "96708412" {
			t.Errorf("Expected height 96708412, got %s", height)
		}
		if organization != "Test Org" {
			t.Errorf("Expected organization 'Test Org', got %s", organization)
		}
		if roleID != "1" {
			t.Errorf("Expected role_id 1, got %s", roleID)
		}
		if sortBy != "tokens_staked" {
			t.Errorf("Expected sort_by 'tokens_staked', got %s", sortBy)
		}

		resp := NodeResponse{
			Data: []Node{
				{NodeID: "abc123", Name: "Test Node"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)

	ctx := context.Background()
	result, err := service.GetNodes().
		Height(96708412).
		Organization("Test Org").
		RoleID("1").
		SortBy("tokens_staked").
		Do(ctx)
	if err != nil {
		t.Fatalf("GetNodes failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 node, got %d", len(result.Data))
	}
}

func TestFlowService_GetNode(t *testing.T) {
	nodeID := "abc123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/node/%s", nodeID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		resp := NodeResponse{
			Data: []Node{
				{
					NodeID:       nodeID,
					Address:      "0x1234",
					Name:         "Test Node",
					Organization: "Test Org",
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
	result, err := service.GetNode().NodeID(nodeID).Do(ctx)
	if err != nil {
		t.Fatalf("GetNode failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 node, got %d", len(result.Data))
	}

	node := result.Data[0]
	if node.NodeID != nodeID {
		t.Errorf("Expected node ID %s, got %s", nodeID, node.NodeID)
	}
}

func TestFlowService_GetNodeDelegationRewards(t *testing.T) {
	nodeID := "abc123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		expectedPath := fmt.Sprintf("/flow/v1/node/%s/reward/delegation", nodeID)
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		address := r.URL.Query().Get("address")
		sortBy := r.URL.Query().Get("sort_by")

		if address != "0x1234" {
			t.Errorf("Expected address 0x1234, got %s", address)
		}
		if sortBy != "amount" {
			t.Errorf("Expected sort_by 'amount', got %s", sortBy)
		}

		resp := DelegationRewardResponse{
			Data: []DelegationReward{
				{
					NodeID:      nodeID,
					Address:     "0x1234",
					Amount:      100.5,
					BlockHeight: 96708412,
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
	result, err := service.GetNodeDelegationRewards().
		NodeID(nodeID).
		Address("0x1234").
		SortBy("amount").
		Do(ctx)
	if err != nil {
		t.Fatalf("GetNodeDelegationRewards failed: %v", err)
	}

	if len(result.Data) != 1 {
		t.Errorf("Expected 1 reward, got %d", len(result.Data))
	}

	reward := result.Data[0]
	if reward.Amount != 100.5 {
		t.Errorf("Expected amount 100.5, got %f", reward.Amount)
	}
	if reward.BlockHeight != 96708412 {
		t.Errorf("Expected block height 96708412, got %d", reward.BlockHeight)
	}
}

func TestFlowService_NodeRequiredFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	defer server.Close()

	client := &mockClient{server: server}
	service := NewService(client)
	ctx := context.Background()

	// Test GetNode without node ID
	_, err := service.GetNode().Do(ctx)
	if err == nil {
		t.Error("Expected error when node ID is not provided")
	}

	// Test GetNodeDelegationRewards without node ID
	_, err = service.GetNodeDelegationRewards().Do(ctx)
	if err == nil {
		t.Error("Expected error when node ID is not provided")
	}
}
