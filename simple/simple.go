package simple

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Client is an interface for making HTTP requests to the API
type Client interface {
	DoRequest(ctx context.Context, method, path string, query url.Values) (*http.Response, error)
	DecodeResponse(resp *http.Response, v any) error
}

// Service handles operations for the Simple API endpoints
type Service struct {
	client Client
}

// NewService creates a new Simple API service
func NewService(client Client) *Service {
	return &Service{client: client}
}

// Block represents a Flow blockchain block
type Block struct {
	Height       uint64          `json:"height"`
	ID           string          `json:"id"`
	Timestamp    string          `json:"timestamp"`
	Transactions []TransactionID `json:"transactions"`
	TxCount      int             `json:"tx"`
}

// TransactionID represents a transaction identifier
type TransactionID struct {
	ID string `json:"id"`
}

// BlocksResponse represents the response from the blocks endpoint
type BlocksResponse struct {
	Blocks []Block `json:"blocks"`
}

// Event represents a Flow blockchain event
type Event struct {
	BlockHeight     uint64                 `json:"block_height"`
	EventIndex      int                    `json:"event_index"`
	Name            string                 `json:"name"`
	Timestamp       string                 `json:"timestamp"`
	TransactionHash string                 `json:"transaction_hash"`
	Fields          map[string]interface{} `json:"fields"`
}

// EventsResponse represents the response from the events endpoint
type EventsResponse struct {
	Events []Event `json:"events"`
}

// Transaction represents a Flow blockchain transaction
type Transaction struct {
	ID                     string                 `json:"id"`
	BlockHeight            uint64                 `json:"block_height"`
	BlockID                string                 `json:"block_id"`
	Timestamp              string                 `json:"timestamp"`
	Payer                  string                 `json:"payer"`
	Proposer               string                 `json:"proposer"`
	ProposerIndex          int                    `json:"proposer_index"`
	ProposerSequenceNumber int                    `json:"proposer_sequence_number"`
	Authorizers            []string               `json:"authorizers"`
	Status                 string                 `json:"status"`
	Error                  string                 `json:"error,omitempty"`
	ErrorCode              string                 `json:"error_code,omitempty"`
	GasLimit               int                    `json:"gas_limit"`
	GasUsed                int                    `json:"gas_used"`
	Fee                    float64                `json:"fee"`
	Argument               interface{}            `json:"argument,omitempty"`
	Events                 []TransactionEvent     `json:"events,omitempty"`
	EventsAggregate        map[string]interface{} `json:"events_aggregate,omitempty"`
	TransactionBody        *TransactionBody       `json:"transaction_body,omitempty"`
}

// TransactionBody contains the transaction script body
type TransactionBody struct {
	Body string `json:"body"`
}

// TransactionEvent represents an event within a transaction
type TransactionEvent struct {
	EventIndex int                    `json:"event_index"`
	Name       string                 `json:"name"`
	Fields     map[string]interface{} `json:"fields"`
}

// TransactionsResponse represents the response from the transaction endpoint
type TransactionsResponse struct {
	Transactions []Transaction `json:"transactions"`
}

// SimpleEvent represents a simplified event structure
type SimpleEvent struct {
	EventIndex int                    `json:"event_index"`
	Name       string                 `json:"name"`
	Fields     map[string]interface{} `json:"fields"`
}

// TransactionEventsResponse represents the response from the transaction events endpoint
type TransactionEventsResponse struct {
	Events []SimpleEvent `json:"events"`
}

// BlocksRequestBuilder builds a request to get blocks
type BlocksRequestBuilder struct {
	service *Service
	height  uint64
	offset  *int
}

// GetBlocks creates a new blocks request builder
func (s *Service) GetBlocks() *BlocksRequestBuilder {
	return &BlocksRequestBuilder{service: s}
}

// Height sets the block height (required)
func (b *BlocksRequestBuilder) Height(height uint64) *BlocksRequestBuilder {
	b.height = height
	return b
}

// Offset sets the pagination offset (optional)
func (b *BlocksRequestBuilder) Offset(offset int) *BlocksRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the blocks request
func (b *BlocksRequestBuilder) Do(ctx context.Context) (*BlocksResponse, error) {
	if b.height == 0 {
		// TODO: we should be able to get the genesis block, but the API currently returns an error
		// {"error":"Field 'Height' failed on the 'required' tag"}
		return nil, fmt.Errorf("height is required")
	}

	query := url.Values{}
	query.Set("height", strconv.FormatUint(b.height, 10))
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/simple/v1/blocks", query)
	if err != nil {
		return nil, err
	}

	var blocksResp BlocksResponse
	if err := b.service.client.DecodeResponse(resp, &blocksResp); err != nil {
		return nil, err
	}

	return &blocksResp, nil
}

// EventsRequestBuilder builds a request to get events
type EventsRequestBuilder struct {
	service    *Service
	name       string
	fromHeight uint64
	toHeight   uint64
	offset     *int
}

// GetEvents creates a new events request builder
func (s *Service) GetEvents() *EventsRequestBuilder {
	return &EventsRequestBuilder{service: s}
}

// Name sets the event name to filter by (required)
func (b *EventsRequestBuilder) Name(name string) *EventsRequestBuilder {
	b.name = name
	return b
}

// FromHeight sets the starting block height (required)
func (b *EventsRequestBuilder) FromHeight(height uint64) *EventsRequestBuilder {
	b.fromHeight = height
	return b
}

// ToHeight sets the ending block height (required)
func (b *EventsRequestBuilder) ToHeight(height uint64) *EventsRequestBuilder {
	b.toHeight = height
	return b
}

// Offset sets the pagination offset (optional)
func (b *EventsRequestBuilder) Offset(offset int) *EventsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the events request
// Returns up to 100 events per request, ordered from oldest to newest
func (b *EventsRequestBuilder) Do(ctx context.Context) (*EventsResponse, error) {
	if b.name == "" {
		return nil, fmt.Errorf("event name is required")
	}
	if b.fromHeight == 0 {
		return nil, fmt.Errorf("from_height is required")
	}
	if b.toHeight == 0 {
		return nil, fmt.Errorf("to_height is required")
	}

	query := url.Values{}
	query.Set("name", b.name)
	query.Set("from_height", strconv.FormatUint(b.fromHeight, 10))
	query.Set("to_height", strconv.FormatUint(b.toHeight, 10))
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/simple/v1/events", query)
	if err != nil {
		return nil, err
	}

	var eventsResp EventsResponse
	if err := b.service.client.DecodeResponse(resp, &eventsResp); err != nil {
		return nil, err
	}

	return &eventsResp, nil
}

// TransactionRequestBuilder builds a request to get a transaction
type TransactionRequestBuilder struct {
	service *Service
	id      string
}

// GetTransaction creates a new transaction request builder
func (s *Service) GetTransaction() *TransactionRequestBuilder {
	return &TransactionRequestBuilder{service: s}
}

// ID sets the transaction ID (required)
func (b *TransactionRequestBuilder) ID(id string) *TransactionRequestBuilder {
	b.id = id
	return b
}

// Do executes the transaction request
func (b *TransactionRequestBuilder) Do(ctx context.Context) (*TransactionsResponse, error) {
	if b.id == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	query := url.Values{}
	query.Set("id", b.id)

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/simple/v1/transaction", query)
	if err != nil {
		return nil, err
	}

	var txResp TransactionsResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}

// TransactionEventsRequestBuilder builds a request to get transaction events
type TransactionEventsRequestBuilder struct {
	service       *Service
	transactionID string
	offset        *int
}

// GetTransactionEvents creates a new transaction events request builder
func (s *Service) GetTransactionEvents() *TransactionEventsRequestBuilder {
	return &TransactionEventsRequestBuilder{service: s}
}

// TransactionID sets the transaction ID (required)
func (b *TransactionEventsRequestBuilder) TransactionID(id string) *TransactionEventsRequestBuilder {
	b.transactionID = id
	return b
}

// Offset sets the pagination offset (optional)
func (b *TransactionEventsRequestBuilder) Offset(offset int) *TransactionEventsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the transaction events request
func (b *TransactionEventsRequestBuilder) Do(ctx context.Context) (*TransactionEventsResponse, error) {
	if b.transactionID == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	query := url.Values{}
	query.Set("transaction_id", b.transactionID)
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/simple/v1/transaction/events", query)
	if err != nil {
		return nil, err
	}

	var eventsResp TransactionEventsResponse
	if err := b.service.client.DecodeResponse(resp, &eventsResp); err != nil {
		return nil, err
	}

	return &eventsResp, nil
}
