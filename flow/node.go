package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Node represents a Flow node
type Node struct {
	Address           string  `json:"address"`
	City              string  `json:"city"`
	Country           string  `json:"country"`
	CountryFlag       string  `json:"country_flag"`
	Delegators        int     `json:"delegators"`
	DelegatorsStaked  float64 `json:"delegators_staked"`
	Epoch             int     `json:"epoch"`
	ID                string  `json:"id"`
	Image             string  `json:"image"`
	IPAddress         string  `json:"ip_address"`
	ISP               string  `json:"isp"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Name              string  `json:"name"`
	NodeID            string  `json:"node_id"`
	Organization      string  `json:"organization"`
	Role              string  `json:"role"`
	RoleID            int     `json:"role_id"`
	TokensStaked      float64 `json:"tokens_staked"`
}

// NodeResponse represents the response from the nodes endpoint
type NodeResponse struct {
	Data  []Node                 `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// DelegationReward represents a delegation reward
type DelegationReward struct {
	Address     string  `json:"address"`
	Amount      float64 `json:"amount"`
	BlockHeight uint64  `json:"block_height"`
	DelegatorID string  `json:"delegator_id"`
	NodeID      string  `json:"node_id"`
	Timestamp   string  `json:"timestamp"`
}

// DelegationRewardResponse represents the response from the delegation rewards endpoint
type DelegationRewardResponse struct {
	Data  []DelegationReward     `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NodesRequestBuilder builds a request to get nodes
type NodesRequestBuilder struct {
	service      *Service
	height       *uint64
	limit        *int
	offset       *int
	organization *string
	roleID       *string
	sortBy       *string
}

// GetNodes creates a new nodes request builder
func (s *Service) GetNodes() *NodesRequestBuilder {
	return &NodesRequestBuilder{service: s}
}

// Height sets the block height filter (optional)
func (b *NodesRequestBuilder) Height(height uint64) *NodesRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 500)
func (b *NodesRequestBuilder) Limit(limit int) *NodesRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *NodesRequestBuilder) Offset(offset int) *NodesRequestBuilder {
	b.offset = &offset
	return b
}

// Organization sets the organization filter (optional)
func (b *NodesRequestBuilder) Organization(organization string) *NodesRequestBuilder {
	b.organization = &organization
	return b
}

// RoleID sets the role ID filter (optional)
// 1 - collection, 2 - consensus, 3 - execution, 4 - verification, 5 - access
func (b *NodesRequestBuilder) RoleID(roleID string) *NodesRequestBuilder {
	b.roleID = &roleID
	return b
}

// SortBy sets the sort field (optional)
// Valid values: 'tokens_staked', 'delegators' (Default = 'block_height')
func (b *NodesRequestBuilder) SortBy(sortBy string) *NodesRequestBuilder {
	b.sortBy = &sortBy
	return b
}

// Do executes the nodes request
func (b *NodesRequestBuilder) Do(ctx context.Context) (*NodeResponse, error) {
	query := url.Values{}
	if b.height != nil {
		query.Set("height", strconv.FormatUint(*b.height, 10))
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.organization != nil {
		query.Set("organiztion", *b.organization) // Note: API has typo "organiztion"
	}
	if b.roleID != nil {
		query.Set("role_id", *b.roleID)
	}
	if b.sortBy != nil {
		query.Set("sort_by", *b.sortBy)
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/node", query)
	if err != nil {
		return nil, err
	}

	var nodeResp NodeResponse
	if err := b.service.client.DecodeResponse(resp, &nodeResp); err != nil {
		return nil, err
	}

	return &nodeResp, nil
}

// NodeRequestBuilder builds a request to get a specific node
type NodeRequestBuilder struct {
	service *Service
	nodeID  string
}

// GetNode creates a new node request builder
func (s *Service) GetNode() *NodeRequestBuilder {
	return &NodeRequestBuilder{service: s}
}

// NodeID sets the node ID (required)
func (b *NodeRequestBuilder) NodeID(nodeID string) *NodeRequestBuilder {
	b.nodeID = nodeID
	return b
}

// Do executes the node request
func (b *NodeRequestBuilder) Do(ctx context.Context) (*NodeResponse, error) {
	if b.nodeID == "" {
		return nil, fmt.Errorf("node ID is required")
	}

	path := fmt.Sprintf("/flow/v1/node/%s", b.nodeID)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var nodeResp NodeResponse
	if err := b.service.client.DecodeResponse(resp, &nodeResp); err != nil {
		return nil, err
	}

	return &nodeResp, nil
}

// NodeDelegationRewardsRequestBuilder builds a request to get delegation rewards for a node
type NodeDelegationRewardsRequestBuilder struct {
	service *Service
	nodeID  string
	limit   *int
	offset  *int
	address *string
	sortBy  *string
}

// GetNodeDelegationRewards creates a new delegation rewards request builder
func (s *Service) GetNodeDelegationRewards() *NodeDelegationRewardsRequestBuilder {
	return &NodeDelegationRewardsRequestBuilder{service: s}
}

// NodeID sets the node ID (required)
func (b *NodeDelegationRewardsRequestBuilder) NodeID(nodeID string) *NodeDelegationRewardsRequestBuilder {
	b.nodeID = nodeID
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *NodeDelegationRewardsRequestBuilder) Limit(limit int) *NodeDelegationRewardsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *NodeDelegationRewardsRequestBuilder) Offset(offset int) *NodeDelegationRewardsRequestBuilder {
	b.offset = &offset
	return b
}

// Address sets the address filter (optional)
func (b *NodeDelegationRewardsRequestBuilder) Address(address string) *NodeDelegationRewardsRequestBuilder {
	b.address = &address
	return b
}

// SortBy sets the sort field (optional)
// Valid values: 'timestamp', 'amount'
func (b *NodeDelegationRewardsRequestBuilder) SortBy(sortBy string) *NodeDelegationRewardsRequestBuilder {
	b.sortBy = &sortBy
	return b
}

// Do executes the delegation rewards request
func (b *NodeDelegationRewardsRequestBuilder) Do(ctx context.Context) (*DelegationRewardResponse, error) {
	if b.nodeID == "" {
		return nil, fmt.Errorf("node ID is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.address != nil {
		query.Set("address", *b.address)
	}
	if b.sortBy != nil {
		query.Set("sort_by", *b.sortBy)
	}

	path := fmt.Sprintf("/flow/v1/node/%s/reward/delegation", b.nodeID)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var rewardResp DelegationRewardResponse
	if err := b.service.client.DecodeResponse(resp, &rewardResp); err != nil {
		return nil, err
	}

	return &rewardResp, nil
}
