package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// FungibleToken represents a fungible token with its details
type FungibleToken struct {
	Address           string  `json:"address"`
	BlockHeight       uint64  `json:"block_height"`
	CirculatingSupply float64 `json:"circulating_supply"`
	Coingecko         string  `json:"coingecko"`
	Coinmarketcap     string  `json:"coinmarketcap"`
	ContractName      string  `json:"contract_name"`
	Decimals          int     `json:"decimals"`
	Description       string  `json:"description"`
	Display           string  `json:"display"`
	External          string  `json:"external"`
	FlowtyID          string  `json:"flowty_id"`
	Icon              string  `json:"icon"`
	IconURL           string  `json:"icon_url"`
	ID                string  `json:"id"`
	Name              string  `json:"name"`
	Path              string  `json:"path"`
	Symbol            string  `json:"symbol"`
	Tags              string  `json:"tags"`
	Token             string  `json:"token"`
	Twitter           string  `json:"twitter"`
	Website           string  `json:"website"`
}

// FTStats represents fungible token statistics
type FTStats struct {
	OwnerCounts  int     `json:"owner_counts"`
	TotalBalance float64 `json:"total_balance"`
}

// FungibleTokenDetails represents detailed information about a fungible token
type FungibleTokenDetails struct {
	FungibleToken
	Stats FTStats `json:"stats"`
}

// FTListResponse represents the response from the fungible tokens list endpoint
type FTListResponse struct {
	Data  []FungibleToken        `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// FungibleTokenResponse represents the response from the fungible token details endpoint
type FungibleTokenResponse struct {
	Data  []FungibleTokenDetails `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// FTTransfer represents a fungible token transfer
type FTTransfer struct {
	Address         string                 `json:"address"`
	Amount          float64                `json:"amount"`
	ApproxUSDPrice  float64                `json:"approx_usd_price"`
	BlockHeight     uint64                 `json:"block_height"`
	Classifier      string                 `json:"classifier"`
	Direction       string                 `json:"direction"`
	IsPrimary       bool                   `json:"is_primary"`
	Receiver        string                 `json:"receiver"`
	ReceiverBalance float64                `json:"receiver_balance"`
	Sender          string                 `json:"sender"`
	Timestamp       string                 `json:"timestamp"`
	Token           FTTransferTokenDetails `json:"token"`
	TransactionHash string                 `json:"transaction_hash"`
	Verified        bool                   `json:"verified"`
}

type FTTransferTokenDetails struct {
	Name   string `json:"name"`
	Token  string `json:"token"`
	Symbol string `json:"symbol"`
	Logo   string `json:"logo"`
}

// TransfersResponse represents the response from the transfers endpoint
type TransfersResponse struct {
	Data  []FTTransfer      `json:"data"`
	Links map[string]string `json:"_links"`
	Meta  map[string]string `json:"_meta,omitempty"`
	Error string            `json:"error,omitempty"`
}

// FTHolding represents a fungible token holding
type FTHolding struct {
	Address    string  `json:"address"`
	Balance    float64 `json:"balance"`
	Percentage float64 `json:"percentage"`
	Token      string  `json:"token"`
}

// FTHoldingResponse represents the response from the holdings endpoint
type FTHoldingResponse struct {
	Data  []FTHolding            `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// Vault represents a token vault for an account
type Vault struct {
	Address     string  `json:"address"`
	Balance     float64 `json:"balance"`
	BlockHeight uint64  `json:"block_height"`
	ID          string  `json:"id"`
	Path        string  `json:"path"`
	Token       string  `json:"token"`
	VaultID     int     `json:"vault_id"`
}

// AccountFungibleTokenResponse represents the response from the account token endpoint
type AccountFungibleTokenResponse struct {
	Data  []Vault                `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// FTsRequestBuilder builds a request to get fungible tokens list
type FTsRequestBuilder struct {
	service *Service
	height  *uint64
	limit   *int
	offset  *int
}

// GetFTs creates a new fungible tokens list request builder
func (s *Service) GetFTs() *FTsRequestBuilder {
	return &FTsRequestBuilder{service: s}
}

// Height sets the block height filter (optional)
func (b *FTsRequestBuilder) Height(height uint64) *FTsRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *FTsRequestBuilder) Limit(limit int) *FTsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *FTsRequestBuilder) Offset(offset int) *FTsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the fungible tokens list request
func (b *FTsRequestBuilder) Do(ctx context.Context) (*FTListResponse, error) {
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

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/ft", query)
	if err != nil {
		return nil, err
	}

	var ftResp FTListResponse
	if err := b.service.client.DecodeResponse(resp, &ftResp); err != nil {
		return nil, err
	}

	return &ftResp, nil
}

// FTRequestBuilder builds a request to get fungible token details
type FTRequestBuilder struct {
	service *Service
	token   string
}

// GetFT creates a new fungible token details request builder
func (s *Service) GetFT() *FTRequestBuilder {
	return &FTRequestBuilder{service: s}
}

// Token sets the token identifier (required)
func (b *FTRequestBuilder) Token(token string) *FTRequestBuilder {
	b.token = token
	return b
}

// Do executes the fungible token details request
func (b *FTRequestBuilder) Do(ctx context.Context) (*FungibleTokenResponse, error) {
	if b.token == "" {
		return nil, fmt.Errorf("token identifier is required")
	}

	path := fmt.Sprintf("/flow/v1/ft/%s", b.token)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var ftResp FungibleTokenResponse
	if err := b.service.client.DecodeResponse(resp, &ftResp); err != nil {
		return nil, err
	}

	return &ftResp, nil
}

// FTTransfersRequestBuilder builds a request to get fungible token transfers
type FTTransfersRequestBuilder struct {
	service         *Service
	token           *string
	transactionHash *string
	height          *uint64
	limit           *int
	offset          *int
}

// GetFTTransfers creates a new fungible token transfers request builder
func (s *Service) GetFTTransfers() *FTTransfersRequestBuilder {
	return &FTTransfersRequestBuilder{service: s}
}

// Token sets the token identifier filter (optional)
func (b *FTTransfersRequestBuilder) Token(token string) *FTTransfersRequestBuilder {
	b.token = &token
	return b
}

// TransactionHash sets the transaction hash filter (optional)
func (b *FTTransfersRequestBuilder) TransactionHash(hash string) *FTTransfersRequestBuilder {
	b.transactionHash = &hash
	return b
}

// Height sets the block height filter (optional)
func (b *FTTransfersRequestBuilder) Height(height uint64) *FTTransfersRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *FTTransfersRequestBuilder) Limit(limit int) *FTTransfersRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *FTTransfersRequestBuilder) Offset(offset int) *FTTransfersRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the fungible token transfers request
func (b *FTTransfersRequestBuilder) Do(ctx context.Context) (*TransfersResponse, error) {
	query := url.Values{}
	if b.token != nil {
		query.Set("token", *b.token)
	}
	if b.transactionHash != nil {
		query.Set("transaction_hash", *b.transactionHash)
	}
	if b.height != nil {
		query.Set("height", strconv.FormatUint(*b.height, 10))
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/ft/transfer", query)
	if err != nil {
		return nil, err
	}

	var transfersResp TransfersResponse
	if err := b.service.client.DecodeResponse(resp, &transfersResp); err != nil {
		return nil, err
	}

	return &transfersResp, nil
}

// FTHoldingsRequestBuilder builds a request to get fungible token holdings
type FTHoldingsRequestBuilder struct {
	service *Service
	token   string
	limit   *int
	offset  *int
}

// GetFTHoldings creates a new fungible token holdings request builder
func (s *Service) GetFTHoldings() *FTHoldingsRequestBuilder {
	return &FTHoldingsRequestBuilder{service: s}
}

// Token sets the token identifier (required)
func (b *FTHoldingsRequestBuilder) Token(token string) *FTHoldingsRequestBuilder {
	b.token = token
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *FTHoldingsRequestBuilder) Limit(limit int) *FTHoldingsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *FTHoldingsRequestBuilder) Offset(offset int) *FTHoldingsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the fungible token holdings request
func (b *FTHoldingsRequestBuilder) Do(ctx context.Context) (*FTHoldingResponse, error) {
	if b.token == "" {
		return nil, fmt.Errorf("token identifier is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/ft/%s/holding", b.token)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var holdingsResp FTHoldingResponse
	if err := b.service.client.DecodeResponse(resp, &holdingsResp); err != nil {
		return nil, err
	}

	return &holdingsResp, nil
}

// FTAccountTokenRequestBuilder builds a request to get account fungible token
type FTAccountTokenRequestBuilder struct {
	service *Service
	token   string
	address string
	limit   *int
	offset  *int
}

// GetFTAccountToken creates a new account fungible token request builder
func (s *Service) GetFTAccountToken() *FTAccountTokenRequestBuilder {
	return &FTAccountTokenRequestBuilder{service: s}
}

// Token sets the token identifier (required)
func (b *FTAccountTokenRequestBuilder) Token(token string) *FTAccountTokenRequestBuilder {
	b.token = token
	return b
}

// Address sets the account address (required)
func (b *FTAccountTokenRequestBuilder) Address(address string) *FTAccountTokenRequestBuilder {
	b.address = address
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *FTAccountTokenRequestBuilder) Limit(limit int) *FTAccountTokenRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *FTAccountTokenRequestBuilder) Offset(offset int) *FTAccountTokenRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account fungible token request
func (b *FTAccountTokenRequestBuilder) Do(ctx context.Context) (*AccountFungibleTokenResponse, error) {
	if b.token == "" {
		return nil, fmt.Errorf("token identifier is required")
	}
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/ft/%s/account/%s", b.token, b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var accountResp AccountFungibleTokenResponse
	if err := b.service.client.DecodeResponse(resp, &accountResp); err != nil {
		return nil, err
	}

	return &accountResp, nil
}
