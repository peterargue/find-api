package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// EvmData represents EVM-related data in a block
type EvmData struct {
	BlockHeight uint64 `json:"block_height"`
}

// Block represents a Flow blockchain block
type Block struct {
	Evm              *EvmData `json:"evm"`
	EvmTxCount       int      `json:"evm_tx_count"`
	Fees             float64  `json:"fees"`
	Height           uint64   `json:"height"`
	ID               string   `json:"id"`
	SurgeFactor      float64  `json:"surge_factor"`
	SystemEventCount int      `json:"system_event_count"`
	Timestamp        string   `json:"timestamp"`
	TotalGasUsed     int      `json:"total_gas_used"`
	Tx               int      `json:"tx"`
}

// BlockResponse represents the response from the blocks list endpoint
type BlockResponse struct {
	Data  []Block                `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// BlockServiceEvent represents a block service event
type BlockServiceEvent struct {
	BlockHeight uint64 `json:"block_height"`
	EventIndex  int    `json:"event_index"`
	EventType   string `json:"event_type"`
	Timestamp   string `json:"timestamp"`
}

// BlockServiceEventResponse represents the response from the block service events endpoint
type BlockServiceEventResponse struct {
	Data  []BlockServiceEvent    `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// BlockTransaction represents a transaction in a block
type BlockTransaction struct {
	Authorizers     []string               `json:"authorizers"`
	BlockHeight     uint64                 `json:"block_height"`
	ContractImports []string               `json:"contract_imports"`
	ContractOutputs []string               `json:"contract_outputs"`
	Entitlements    []string               `json:"entitlements"`
	Error           string                 `json:"error"`
	ErrorCode       string                 `json:"error_code"`
	EventCount      int                    `json:"event_count"`
	Events          []interface{}          `json:"events,omitempty"`
	EvmTxCount      int                    `json:"evm_tx_count"`
	Fee             float64                `json:"fee"`
	GasLimit        int                    `json:"gas_limit"`
	GasUsed         int                    `json:"gas_used"`
	Payer           string                 `json:"payer"`
	Proposer        string                 `json:"proposer"`
	Roles           map[string]interface{} `json:"roles"`
	Status          string                 `json:"status"`
	Timestamp       string                 `json:"timestamp"`
	TransactionHash string                 `json:"transaction_hash"`
	TransactionID   string                 `json:"transaction_id"`
}

// BlockTransactionsResponse represents the response from the block transactions endpoint
type BlockTransactionsResponse struct {
	Data  []BlockTransaction     `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// BlocksRequestBuilder builds a request to get blocks list
type BlocksRequestBuilder struct {
	service *Service
	height  *uint64
	limit   *int
	offset  *int
}

// GetBlocks creates a new blocks list request builder
func (s *Service) GetBlocks() *BlocksRequestBuilder {
	return &BlocksRequestBuilder{service: s}
}

// Height sets the block height to start from (optional, descending)
func (b *BlocksRequestBuilder) Height(height uint64) *BlocksRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *BlocksRequestBuilder) Limit(limit int) *BlocksRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *BlocksRequestBuilder) Offset(offset int) *BlocksRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the blocks list request
func (b *BlocksRequestBuilder) Do(ctx context.Context) (*BlockResponse, error) {
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

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/block", query)
	if err != nil {
		return nil, err
	}

	var blockResp BlockResponse
	if err := b.service.client.DecodeResponse(resp, &blockResp); err != nil {
		return nil, err
	}

	return &blockResp, nil
}

// BlockRequestBuilder builds a request to get a specific block by height
type BlockRequestBuilder struct {
	service *Service
	height  uint64
}

// GetBlock creates a new block request builder
func (s *Service) GetBlock() *BlockRequestBuilder {
	return &BlockRequestBuilder{service: s}
}

// Height sets the block height (required)
func (b *BlockRequestBuilder) Height(height uint64) *BlockRequestBuilder {
	b.height = height
	return b
}

// Do executes the block request
func (b *BlockRequestBuilder) Do(ctx context.Context) (*BlockResponse, error) {
	if b.height == 0 {
		return nil, fmt.Errorf("block height is required")
	}

	path := fmt.Sprintf("/flow/v1/block/%d", b.height)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var blockResp BlockResponse
	if err := b.service.client.DecodeResponse(resp, &blockResp); err != nil {
		return nil, err
	}

	return &blockResp, nil
}

// BlockServiceEventsRequestBuilder builds a request to get block service events
type BlockServiceEventsRequestBuilder struct {
	service *Service
	height  uint64
	limit   *int
	offset  *int
}

// GetBlockServiceEvents creates a new block service events request builder
func (s *Service) GetBlockServiceEvents() *BlockServiceEventsRequestBuilder {
	return &BlockServiceEventsRequestBuilder{service: s}
}

// Height sets the block height (required)
func (b *BlockServiceEventsRequestBuilder) Height(height uint64) *BlockServiceEventsRequestBuilder {
	b.height = height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *BlockServiceEventsRequestBuilder) Limit(limit int) *BlockServiceEventsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *BlockServiceEventsRequestBuilder) Offset(offset int) *BlockServiceEventsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the block service events request
func (b *BlockServiceEventsRequestBuilder) Do(ctx context.Context) (*BlockServiceEventResponse, error) {
	if b.height == 0 {
		return nil, fmt.Errorf("block height is required")
	}

	query := url.Values{}
	query.Set("height", strconv.FormatUint(b.height, 10))
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/block/%d/service-event", b.height)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var eventsResp BlockServiceEventResponse
	if err := b.service.client.DecodeResponse(resp, &eventsResp); err != nil {
		return nil, err
	}

	return &eventsResp, nil
}

// BlockTransactionsRequestBuilder builds a request to get block transactions
type BlockTransactionsRequestBuilder struct {
	service       *Service
	height        uint64
	includeEvents *bool
}

// GetBlockTransactions creates a new block transactions request builder
func (s *Service) GetBlockTransactions() *BlockTransactionsRequestBuilder {
	return &BlockTransactionsRequestBuilder{service: s}
}

// Height sets the block height (required)
func (b *BlockTransactionsRequestBuilder) Height(height uint64) *BlockTransactionsRequestBuilder {
	b.height = height
	return b
}

// IncludeEvents sets whether to include events in the response (optional, default false)
func (b *BlockTransactionsRequestBuilder) IncludeEvents(include bool) *BlockTransactionsRequestBuilder {
	b.includeEvents = &include
	return b
}

// Do executes the block transactions request
func (b *BlockTransactionsRequestBuilder) Do(ctx context.Context) (*BlockTransactionsResponse, error) {
	if b.height == 0 {
		return nil, fmt.Errorf("block height is required")
	}

	query := url.Values{}
	if b.includeEvents != nil {
		query.Set("include_events", strconv.FormatBool(*b.includeEvents))
	}

	path := fmt.Sprintf("/flow/v1/block/%d/transaction", b.height)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var txResp BlockTransactionsResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}
