package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Contract represents a contract
type Contract struct {
	Address      string `json:"address"`
	BlockHeight  uint64 `json:"block_height"`
	ContractName string `json:"contract_name"`
	Identifier   string `json:"identifier"`
}

// ContractResponse represents the response from the contracts endpoint
type ContractResponse struct {
	Data  []Contract             `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// ContractsRequestBuilder builds a request to get contracts
type ContractsRequestBuilder struct {
	service *Service
	limit   *int
	offset  *int
}

// GetContracts creates a new contracts request builder
func (s *Service) GetContracts() *ContractsRequestBuilder {
	return &ContractsRequestBuilder{service: s}
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *ContractsRequestBuilder) Limit(limit int) *ContractsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *ContractsRequestBuilder) Offset(offset int) *ContractsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the contracts request
func (b *ContractsRequestBuilder) Do(ctx context.Context) (*ContractResponse, error) {
	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/contract", query)
	if err != nil {
		return nil, err
	}

	var contractResp ContractResponse
	if err := b.service.client.DecodeResponse(resp, &contractResp); err != nil {
		return nil, err
	}

	return &contractResp, nil
}

// ContractsByIdentifierRequestBuilder builds a request to get contracts by identifier
type ContractsByIdentifierRequestBuilder struct {
	service    *Service
	identifier string
	limit      *int
	offset     *int
}

// GetContractsByIdentifier creates a new contracts by identifier request builder
func (s *Service) GetContractsByIdentifier() *ContractsByIdentifierRequestBuilder {
	return &ContractsByIdentifierRequestBuilder{service: s}
}

// Identifier sets the contract identifier (required)
func (b *ContractsByIdentifierRequestBuilder) Identifier(identifier string) *ContractsByIdentifierRequestBuilder {
	b.identifier = identifier
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *ContractsByIdentifierRequestBuilder) Limit(limit int) *ContractsByIdentifierRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *ContractsByIdentifierRequestBuilder) Offset(offset int) *ContractsByIdentifierRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the contracts by identifier request
func (b *ContractsByIdentifierRequestBuilder) Do(ctx context.Context) (*ContractResponse, error) {
	if b.identifier == "" {
		return nil, fmt.Errorf("contract identifier is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/contract/%s", b.identifier)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var contractResp ContractResponse
	if err := b.service.client.DecodeResponse(resp, &contractResp); err != nil {
		return nil, err
	}

	return &contractResp, nil
}

// ContractRequestBuilder builds a request to get a specific contract
type ContractRequestBuilder struct {
	service    *Service
	identifier string
	id         string
}

// GetContract creates a new contract request builder
func (s *Service) GetContract() *ContractRequestBuilder {
	return &ContractRequestBuilder{service: s}
}

// Identifier sets the contract identifier (required)
func (b *ContractRequestBuilder) Identifier(identifier string) *ContractRequestBuilder {
	b.identifier = identifier
	return b
}

// ID sets the contract ID (required)
func (b *ContractRequestBuilder) ID(id string) *ContractRequestBuilder {
	b.id = id
	return b
}

// Do executes the contract request
func (b *ContractRequestBuilder) Do(ctx context.Context) (*ContractResponse, error) {
	if b.identifier == "" {
		return nil, fmt.Errorf("contract identifier is required")
	}
	if b.id == "" {
		return nil, fmt.Errorf("contract ID is required")
	}

	path := fmt.Sprintf("/flow/v1/contract/%s/%s", b.identifier, b.id)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var contractResp ContractResponse
	if err := b.service.client.DecodeResponse(resp, &contractResp); err != nil {
		return nil, err
	}

	return &contractResp, nil
}
