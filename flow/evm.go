package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// EvmToken represents an EVM token
type EvmToken struct {
	ContractAddressHash string `json:"contract_address_hash"`
	Decimals            int    `json:"decimals"`
	Holders             int    `json:"holders"`
	IconURL             string `json:"icon_url"`
	Name                string `json:"name"`
	Symbol              string `json:"symbol"`
	TotalSupply         string `json:"total_supply"`
	Transfers           int    `json:"transfers"`
	Type                string `json:"type"`
}

// EvmTokenResponse represents the response from the EVM tokens endpoint
type EvmTokenResponse struct {
	Data  []EvmToken             `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// EvmTransaction represents an EVM transaction
type EvmTransaction struct {
	BlockNumber                     uint64 `json:"block_number"`
	From                            string `json:"from"`
	GasLimit                        string `json:"gas_limit"`
	GasPrice                        string `json:"gas_price"`
	GasUsed                         string `json:"gas_used"`
	HasErrorInInternalTransactions  bool   `json:"has_error_in_internal_transactions"`
	Hash                            string `json:"hash"`
	MaxFeePerGas                    string `json:"max_fee_per_gas"`
	MaxPriorityFeePerGas            string `json:"max_priority_fee_per_gas"`
	Nonce                           int    `json:"nonce"`
	R                               string `json:"r"`
	S                               string `json:"s"`
	Status                          string `json:"status"`
	Timestamp                       string `json:"timestamp"`
	To                              string `json:"to"`
	TransactionIndex                int    `json:"transaction_index"`
	Type                            int    `json:"type"`
	V                               string `json:"v"`
	Value                           string `json:"value"`
}

// EvmTransactionResponse represents the response from the EVM transactions list endpoint
type EvmTransactionResponse struct {
	Data  []EvmTransaction       `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// EvmTokensRequestBuilder builds a request to get EVM tokens
type EvmTokensRequestBuilder struct {
	service *Service
	typ     *string
	name    *string
	limit   *int
	offset  *int
}

// GetEvmTokens creates a new EVM tokens request builder
func (s *Service) GetEvmTokens() *EvmTokensRequestBuilder {
	return &EvmTokensRequestBuilder{service: s}
}

// Type sets the token type filter (optional)
func (b *EvmTokensRequestBuilder) Type(typ string) *EvmTokensRequestBuilder {
	b.typ = &typ
	return b
}

// Name sets the partial name or symbol to search for (optional)
func (b *EvmTokensRequestBuilder) Name(name string) *EvmTokensRequestBuilder {
	b.name = &name
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *EvmTokensRequestBuilder) Limit(limit int) *EvmTokensRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *EvmTokensRequestBuilder) Offset(offset int) *EvmTokensRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the EVM tokens request
func (b *EvmTokensRequestBuilder) Do(ctx context.Context) (*EvmTokenResponse, error) {
	query := url.Values{}
	if b.typ != nil {
		query.Set("type", *b.typ)
	}
	if b.name != nil {
		query.Set("name", *b.name)
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/evm/token", query)
	if err != nil {
		return nil, err
	}

	var tokenResp EvmTokenResponse
	if err := b.service.client.DecodeResponse(resp, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// EvmTokenRequestBuilder builds a request to get a specific EVM token by address
type EvmTokenRequestBuilder struct {
	service *Service
	address string
	limit   *int
	offset  *int
}

// GetEvmToken creates a new EVM token request builder
func (s *Service) GetEvmToken() *EvmTokenRequestBuilder {
	return &EvmTokenRequestBuilder{service: s}
}

// Address sets the token contract address (required)
func (b *EvmTokenRequestBuilder) Address(address string) *EvmTokenRequestBuilder {
	b.address = address
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *EvmTokenRequestBuilder) Limit(limit int) *EvmTokenRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *EvmTokenRequestBuilder) Offset(offset int) *EvmTokenRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the EVM token request
func (b *EvmTokenRequestBuilder) Do(ctx context.Context) (*EvmTokenResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("token address is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/evm/token/%s", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var tokenResp EvmTokenResponse
	if err := b.service.client.DecodeResponse(resp, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// EvmTransactionsRequestBuilder builds a request to get EVM transactions
type EvmTransactionsRequestBuilder struct {
	service *Service
	height  *uint64
	limit   *int
	offset  *int
}

// GetEvmTransactions creates a new EVM transactions request builder
func (s *Service) GetEvmTransactions() *EvmTransactionsRequestBuilder {
	return &EvmTransactionsRequestBuilder{service: s}
}

// Height sets the block height filter (optional)
func (b *EvmTransactionsRequestBuilder) Height(height uint64) *EvmTransactionsRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *EvmTransactionsRequestBuilder) Limit(limit int) *EvmTransactionsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional, requires height parameter)
func (b *EvmTransactionsRequestBuilder) Offset(offset int) *EvmTransactionsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the EVM transactions request
func (b *EvmTransactionsRequestBuilder) Do(ctx context.Context) (*EvmTransactionResponse, error) {
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

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/evm/transaction", query)
	if err != nil {
		return nil, err
	}

	var txResp EvmTransactionResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}

// EvmTransactionRequestBuilder builds a request to get a specific EVM transaction by hash
type EvmTransactionRequestBuilder struct {
	service *Service
	hash    string
}

// GetEvmTransaction creates a new EVM transaction request builder
func (s *Service) GetEvmTransaction() *EvmTransactionRequestBuilder {
	return &EvmTransactionRequestBuilder{service: s}
}

// Hash sets the transaction hash (required)
func (b *EvmTransactionRequestBuilder) Hash(hash string) *EvmTransactionRequestBuilder {
	b.hash = hash
	return b
}

// Do executes the EVM transaction request
func (b *EvmTransactionRequestBuilder) Do(ctx context.Context) (*EvmTransaction, error) {
	if b.hash == "" {
		return nil, fmt.Errorf("transaction hash is required")
	}

	path := fmt.Sprintf("/flow/v1/evm/transaction/%s", b.hash)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var tx EvmTransaction
	if err := b.service.client.DecodeResponse(resp, &tx); err != nil {
		return nil, err
	}

	return &tx, nil
}
