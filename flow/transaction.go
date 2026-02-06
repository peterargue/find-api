package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Transaction represents a Flow transaction in list format
type Transaction struct {
	Authorizers      []string  `json:"authorizers"`
	BlockHeight      uint64    `json:"block_height"`
	ContractImports  []string  `json:"contract_imports"`
	ContractOutputs  []string  `json:"contract_outputs"`
	Error            string    `json:"error"`
	ErrorCode        string    `json:"error_code"`
	EventCount       int       `json:"event_count"`
	Events           []Event   `json:"events,omitempty"`
	Fee              float64   `json:"fee"`
	GasUsed          int       `json:"gas_used"`
	ID               string    `json:"id"`
	Payer            string    `json:"payer"`
	Proposer         string    `json:"proposer"`
	Status           string    `json:"status"`
	SurgeFactor      float64   `json:"surge_factor"`
	Tags             []Tag     `json:"tags"`
	Timestamp        string    `json:"timestamp"`
	TransactionIndex int       `json:"transaction_index"`
	Type             string    `json:"type"`
}

// Event represents a transaction event
type Event struct {
	BlockHeight uint64      `json:"block_height"`
	EventIndex  int         `json:"event_index"`
	Fields      interface{} `json:"fields"`
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Timestamp   string      `json:"timestamp"`
}

// Tag represents a transaction tag
type Tag struct {
	ID   string `json:"id"`
	Logo string `json:"logo"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// TransactionDetails represents detailed transaction information
type TransactionDetails struct {
	Argument         []ArgumentItem     `json:"argument"`
	Authorizers      []string           `json:"authorizers"`
	BlockHeight      uint64             `json:"block_height"`
	BlockID          string             `json:"block_id"`
	ContractImports  []string           `json:"contract_imports"`
	ContractOutputs  []string           `json:"contract_outputs"`
	Error            string             `json:"error"`
	ErrorCode        string             `json:"error_code"`
	Events           []EventOutput      `json:"events"`
	EvmTransactions  []EvmTransactions  `json:"evm_transactions"`
	ExecutionEffort  float64            `json:"execution_effort"`
	Fee              float64            `json:"fee"`
	GasUsed          int                `json:"gas_used"`
	ID               string             `json:"id"`
	Imports          []ImportOutput     `json:"imports"`
	Payer            string             `json:"payer"`
	Proposer         string             `json:"proposer"`
	Script           string             `json:"script"`
	Status           string             `json:"status"`
	SurgeFactor      float64            `json:"surge_factor"`
	Tags             []Tag              `json:"tags"`
	Timestamp        string             `json:"timestamp"`
}

// ArgumentItem represents a transaction argument
type ArgumentItem struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

// EventOutput represents an event in transaction details
type EventOutput struct {
	Data             map[string]interface{} `json:"data"`
	EventIndex       int                    `json:"event_index"`
	TransactionID    string                 `json:"transaction_id"`
	TransactionIndex int                    `json:"transaction_index"`
	Type             string                 `json:"type"`
}

// EvmTransactions represents EVM transaction information
type EvmTransactions struct {
	BlockNumber uint64 `json:"block_number"`
	Hash        string `json:"hash"`
}

// ImportOutput represents contract import information
type ImportOutput struct {
	Location string `json:"location"`
	Name     string `json:"name"`
}

// TransactionsResponse represents the response from the transactions list endpoint
type TransactionsResponse struct {
	Data  []Transaction          `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// TransactionResponse represents the response from the transaction details endpoint
type TransactionResponse struct {
	Data  []TransactionDetails   `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// TransactionsRequestBuilder builds a request to get transactions
type TransactionsRequestBuilder struct {
	service            *Service
	authorizers        *string
	contractIdentifier *string
	from               *string
	height             *uint64
	includeEvents      *bool
	limit              *int
	maxEvents          *int
	maxGas             *int
	minEvents          *int
	minGas             *int
	offset             *int
	payer              *string
	proposer           *string
	status             *string
	to                 *string
	typ                *string
}

// GetTransactions creates a new transactions request builder
func (s *Service) GetTransactions() *TransactionsRequestBuilder {
	return &TransactionsRequestBuilder{service: s}
}

// Authorizers sets the authorizer address filter (optional)
func (b *TransactionsRequestBuilder) Authorizers(authorizers string) *TransactionsRequestBuilder {
	b.authorizers = &authorizers
	return b
}

// ContractIdentifier sets the contract identifier filter (optional)
func (b *TransactionsRequestBuilder) ContractIdentifier(contractIdentifier string) *TransactionsRequestBuilder {
	b.contractIdentifier = &contractIdentifier
	return b
}

// From sets the start timestamp filter (optional, ISO 8601 format)
func (b *TransactionsRequestBuilder) From(from string) *TransactionsRequestBuilder {
	b.from = &from
	return b
}

// Height sets the block height filter (optional)
func (b *TransactionsRequestBuilder) Height(height uint64) *TransactionsRequestBuilder {
	b.height = &height
	return b
}

// IncludeEvents sets whether to include events in the response (optional, default false)
func (b *TransactionsRequestBuilder) IncludeEvents(includeEvents bool) *TransactionsRequestBuilder {
	b.includeEvents = &includeEvents
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *TransactionsRequestBuilder) Limit(limit int) *TransactionsRequestBuilder {
	b.limit = &limit
	return b
}

// MaxEvents sets the maximum number of events filter (optional)
func (b *TransactionsRequestBuilder) MaxEvents(maxEvents int) *TransactionsRequestBuilder {
	b.maxEvents = &maxEvents
	return b
}

// MaxGas sets the maximum gas used filter (optional)
func (b *TransactionsRequestBuilder) MaxGas(maxGas int) *TransactionsRequestBuilder {
	b.maxGas = &maxGas
	return b
}

// MinEvents sets the minimum number of events filter (optional)
func (b *TransactionsRequestBuilder) MinEvents(minEvents int) *TransactionsRequestBuilder {
	b.minEvents = &minEvents
	return b
}

// MinGas sets the minimum gas used filter (optional)
func (b *TransactionsRequestBuilder) MinGas(minGas int) *TransactionsRequestBuilder {
	b.minGas = &minGas
	return b
}

// Offset sets the pagination offset (optional)
func (b *TransactionsRequestBuilder) Offset(offset int) *TransactionsRequestBuilder {
	b.offset = &offset
	return b
}

// Payer sets the payer address filter (optional)
func (b *TransactionsRequestBuilder) Payer(payer string) *TransactionsRequestBuilder {
	b.payer = &payer
	return b
}

// Proposer sets the proposer address filter (optional)
func (b *TransactionsRequestBuilder) Proposer(proposer string) *TransactionsRequestBuilder {
	b.proposer = &proposer
	return b
}

// Status sets the status filter (optional, e.g., ERROR, SEALED)
func (b *TransactionsRequestBuilder) Status(status string) *TransactionsRequestBuilder {
	b.status = &status
	return b
}

// To sets the end timestamp filter (optional, ISO 8601 format)
func (b *TransactionsRequestBuilder) To(to string) *TransactionsRequestBuilder {
	b.to = &to
	return b
}

// Type sets the transaction type filter (optional)
func (b *TransactionsRequestBuilder) Type(typ string) *TransactionsRequestBuilder {
	b.typ = &typ
	return b
}

// Do executes the transactions request
func (b *TransactionsRequestBuilder) Do(ctx context.Context) (*TransactionsResponse, error) {
	query := url.Values{}
	if b.authorizers != nil {
		query.Set("authorizers", *b.authorizers)
	}
	if b.contractIdentifier != nil {
		query.Set("contract_identifier", *b.contractIdentifier)
	}
	if b.from != nil {
		query.Set("from", *b.from)
	}
	if b.height != nil {
		query.Set("height", strconv.FormatUint(*b.height, 10))
	}
	if b.includeEvents != nil {
		query.Set("include_events", strconv.FormatBool(*b.includeEvents))
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.maxEvents != nil {
		query.Set("max_events", strconv.Itoa(*b.maxEvents))
	}
	if b.maxGas != nil {
		query.Set("max_gas", strconv.Itoa(*b.maxGas))
	}
	if b.minEvents != nil {
		query.Set("min_events", strconv.Itoa(*b.minEvents))
	}
	if b.minGas != nil {
		query.Set("min_gas", strconv.Itoa(*b.minGas))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.payer != nil {
		query.Set("payer", *b.payer)
	}
	if b.proposer != nil {
		query.Set("proposer", *b.proposer)
	}
	if b.status != nil {
		query.Set("status", *b.status)
	}
	if b.to != nil {
		query.Set("to", *b.to)
	}
	if b.typ != nil {
		query.Set("type", *b.typ)
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/transaction", query)
	if err != nil {
		return nil, err
	}

	var txResp TransactionsResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}

// TransactionRequestBuilder builds a request to get a specific transaction
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
func (b *TransactionRequestBuilder) Do(ctx context.Context) (*TransactionResponse, error) {
	if b.id == "" {
		return nil, fmt.Errorf("transaction ID is required")
	}

	path := fmt.Sprintf("/flow/v1/transaction/%s", b.id)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var txResp TransactionResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}

// ScheduledTransaction represents a scheduled transaction
type ScheduledTransaction struct {
	Args                   map[string]interface{} `json:"args"`
	BlockHeight            uint64                 `json:"block_height"`
	CompletedAt            string                 `json:"completed_at"`
	CompletedBlockHeight   uint64                 `json:"completed_block_height"`
	CompletedTransaction   string                 `json:"completed_transaction"`
	CreatedAt              string                 `json:"created_at"`
	Error                  string                 `json:"error"`
	ExecutionEffort        int                    `json:"execution_effort"`
	Fees                   string                 `json:"fees"`
	Handler                string                 `json:"handler"`
	HandlerContract        string                 `json:"handler_contract"`
	HandlerUUID            int                    `json:"handler_uuid"`
	ID                     string                 `json:"id"`
	IsCompleted            bool                   `json:"is_completed"`
	Owner                  string                 `json:"owner"`
	PinChanged             bool                   `json:"pin_changed"`
	Priority               string                 `json:"priority"`
	ScheduledAt            string                 `json:"scheduled_at"`
	ScheduledTransaction   string                 `json:"scheduled_transaction"`
	Status                 string                 `json:"status"`
}

// ScheduledTransactionsResponse represents the response from the scheduled transactions endpoint
type ScheduledTransactionsResponse struct {
	Data  []ScheduledTransaction `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// ScheduledTransactionsRequestBuilder builds a request to get scheduled transactions
type ScheduledTransactionsRequestBuilder struct {
	service              *Service
	completed            *bool
	completedFrom        *string
	completedHeight      *uint64
	completedTo          *string
	contractIdentifier   *string
	handler              *string
	handlerUUID          *int
	heightFrom           *uint64
	heightTo             *uint64
	id                   *string
	limit                *int
	offset               *int
	owner                *string
	priority             *string
	scheduledFrom        *string
	scheduledTo          *string
	status               *string
}

// GetScheduledTransactions creates a new scheduled transactions request builder
func (s *Service) GetScheduledTransactions() *ScheduledTransactionsRequestBuilder {
	return &ScheduledTransactionsRequestBuilder{service: s}
}

// Completed sets the completed filter (optional)
func (b *ScheduledTransactionsRequestBuilder) Completed(completed bool) *ScheduledTransactionsRequestBuilder {
	b.completed = &completed
	return b
}

// CompletedFrom sets the completed from timestamp filter (optional)
func (b *ScheduledTransactionsRequestBuilder) CompletedFrom(completedFrom string) *ScheduledTransactionsRequestBuilder {
	b.completedFrom = &completedFrom
	return b
}

// CompletedHeight sets the completed height filter (optional)
func (b *ScheduledTransactionsRequestBuilder) CompletedHeight(completedHeight uint64) *ScheduledTransactionsRequestBuilder {
	b.completedHeight = &completedHeight
	return b
}

// CompletedTo sets the completed to timestamp filter (optional)
func (b *ScheduledTransactionsRequestBuilder) CompletedTo(completedTo string) *ScheduledTransactionsRequestBuilder {
	b.completedTo = &completedTo
	return b
}

// ContractIdentifier sets the contract identifier filter (optional)
func (b *ScheduledTransactionsRequestBuilder) ContractIdentifier(contractIdentifier string) *ScheduledTransactionsRequestBuilder {
	b.contractIdentifier = &contractIdentifier
	return b
}

// Handler sets the handler filter (optional)
func (b *ScheduledTransactionsRequestBuilder) Handler(handler string) *ScheduledTransactionsRequestBuilder {
	b.handler = &handler
	return b
}

// HandlerUUID sets the handler UUID filter (optional)
func (b *ScheduledTransactionsRequestBuilder) HandlerUUID(handlerUUID int) *ScheduledTransactionsRequestBuilder {
	b.handlerUUID = &handlerUUID
	return b
}

// HeightFrom sets the height from filter (optional)
func (b *ScheduledTransactionsRequestBuilder) HeightFrom(heightFrom uint64) *ScheduledTransactionsRequestBuilder {
	b.heightFrom = &heightFrom
	return b
}

// HeightTo sets the height to filter (optional)
func (b *ScheduledTransactionsRequestBuilder) HeightTo(heightTo uint64) *ScheduledTransactionsRequestBuilder {
	b.heightTo = &heightTo
	return b
}

// ID sets the transaction ID filter (optional)
func (b *ScheduledTransactionsRequestBuilder) ID(id string) *ScheduledTransactionsRequestBuilder {
	b.id = &id
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *ScheduledTransactionsRequestBuilder) Limit(limit int) *ScheduledTransactionsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *ScheduledTransactionsRequestBuilder) Offset(offset int) *ScheduledTransactionsRequestBuilder {
	b.offset = &offset
	return b
}

// Owner sets the owner filter (optional)
func (b *ScheduledTransactionsRequestBuilder) Owner(owner string) *ScheduledTransactionsRequestBuilder {
	b.owner = &owner
	return b
}

// Priority sets the priority filter (optional)
// Valid values: low, medium, high
func (b *ScheduledTransactionsRequestBuilder) Priority(priority string) *ScheduledTransactionsRequestBuilder {
	b.priority = &priority
	return b
}

// ScheduledFrom sets the scheduled from timestamp filter (optional)
func (b *ScheduledTransactionsRequestBuilder) ScheduledFrom(scheduledFrom string) *ScheduledTransactionsRequestBuilder {
	b.scheduledFrom = &scheduledFrom
	return b
}

// ScheduledTo sets the scheduled to timestamp filter (optional)
func (b *ScheduledTransactionsRequestBuilder) ScheduledTo(scheduledTo string) *ScheduledTransactionsRequestBuilder {
	b.scheduledTo = &scheduledTo
	return b
}

// Status sets the status filter (optional)
func (b *ScheduledTransactionsRequestBuilder) Status(status string) *ScheduledTransactionsRequestBuilder {
	b.status = &status
	return b
}

// Do executes the scheduled transactions request
func (b *ScheduledTransactionsRequestBuilder) Do(ctx context.Context) (*ScheduledTransactionsResponse, error) {
	query := url.Values{}
	if b.completed != nil {
		query.Set("completed", strconv.FormatBool(*b.completed))
	}
	if b.completedFrom != nil {
		query.Set("completed_from", *b.completedFrom)
	}
	if b.completedHeight != nil {
		query.Set("completed_height", strconv.FormatUint(*b.completedHeight, 10))
	}
	if b.completedTo != nil {
		query.Set("completed_to", *b.completedTo)
	}
	if b.contractIdentifier != nil {
		query.Set("contract_identifier", *b.contractIdentifier)
	}
	if b.handler != nil {
		query.Set("handler", *b.handler)
	}
	if b.handlerUUID != nil {
		query.Set("handler_uuid", strconv.Itoa(*b.handlerUUID))
	}
	if b.heightFrom != nil {
		query.Set("height_from", strconv.FormatUint(*b.heightFrom, 10))
	}
	if b.heightTo != nil {
		query.Set("height_to", strconv.FormatUint(*b.heightTo, 10))
	}
	if b.id != nil {
		query.Set("id", *b.id)
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.owner != nil {
		query.Set("owner", *b.owner)
	}
	if b.priority != nil {
		query.Set("priority", *b.priority)
	}
	if b.scheduledFrom != nil {
		query.Set("scheduled_from", *b.scheduledFrom)
	}
	if b.scheduledTo != nil {
		query.Set("scheduled_to", *b.scheduledTo)
	}
	if b.status != nil {
		query.Set("status", *b.status)
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/scheduled-transaction", query)
	if err != nil {
		return nil, err
	}

	var scheduledResp ScheduledTransactionsResponse
	if err := b.service.client.DecodeResponse(resp, &scheduledResp); err != nil {
		return nil, err
	}

	return &scheduledResp, nil
}
