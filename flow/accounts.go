package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// Account represents basic account information
type Account struct {
	Address          string                 `json:"address"`
	Creator          string                 `json:"creator"`
	Data             map[string]interface{} `json:"data"`
	FindName         string                 `json:"find_name"`
	FlowBalance      float64                `json:"flow_balance"`
	FlowStorage      float64                `json:"flow_storage"`
	Height           uint64                 `json:"height"`
	StorageAvailable float64                `json:"storage_available"`
	StorageUsed      float64                `json:"storage_used"`
	Timestamp        string                 `json:"timestamp"`
	TransactionHash  string                 `json:"transaction_hash"`
}

// AccountsResponse represents the response from the accounts list endpoint
type AccountsResponse struct {
	Data  []Account              `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// AccountInfo represents on-chain account information
type AccountInfo struct {
	DelegatedBalance float64 `json:"delegatedBalance"`
	LockedBalance    float64 `json:"lockedBalance"`
	StakedBalance    float64 `json:"stakedBalance"`
	UnlockedBalance  float64 `json:"unlockedBalance"`
}

// VaultInfo represents vault information
type VaultInfo struct {
	Balance float64 `json:"balance"`
	Path    string  `json:"path"`
}

// KeyInfo represents account key information
type KeyInfo struct {
	Index          int    `json:"index"`
	PublicKey      string `json:"publicKey"`
	SignAlgo       string `json:"signAlgo"`
	HashAlgo       string `json:"hashAlgo"`
	Weight         int    `json:"weight"`
	SequenceNumber int    `json:"sequenceNumber"`
	Revoked        bool   `json:"revoked"`
}

// Domains represents domain information
type Domains struct {
	FindName string `json:"findName"`
}

// Find represents find name information
type Find struct {
	Name string `json:"name"`
}

// CombinedAccountDetails represents detailed account information
type CombinedAccountDetails struct {
	AccountInfo      *AccountInfo           `json:"accountInfo"`
	Address          string                 `json:"address"`
	Contracts        []string               `json:"contracts"`
	Domains          *Domains               `json:"domains"`
	Find             *Find                  `json:"find"`
	FlowBalance      float64                `json:"flowBalance"`
	FlowStorage      float64                `json:"flowStorage"`
	Keys             []KeyInfo              `json:"keys"`
	StorageAvailable float64                `json:"storageAvailable"`
	StorageUsed      float64                `json:"storageUsed"`
	Vaults           map[string]VaultInfo   `json:"vaults"`
}

// AccountDetailsResponse represents the response from the account details endpoint
type AccountDetailsResponse struct {
	Data  []CombinedAccountDetails `json:"data"`
	Links map[string]string        `json:"_links"`
	Meta  map[string]interface{}   `json:"_meta"`
	Error interface{}              `json:"error,omitempty"`
}

// AccountFTCollection represents an FT collection in an account
type AccountFTCollection struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Path    string `json:"path"`
	Token   string `json:"token"`
	VaultID int    `json:"vault_id"`
}

// AccountFTCollectionsResponse represents the response from the account FT collections endpoint
type AccountFTCollectionsResponse struct {
	Data  []AccountFTCollection  `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// AccountTransaction represents a transaction for an account
type AccountTransaction struct {
	Authorizers     []string               `json:"authorizers"`
	BlockHeight     uint64                 `json:"block_height"`
	ContractImports []string               `json:"contract_imports"`
	ContractOutputs []string               `json:"contract_outputs"`
	Entitlements    []string               `json:"entitlements"`
	Error           string                 `json:"error"`
	ErrorCode       string                 `json:"error_code"`
	EventCount      int                    `json:"event_count"`
	Events          []interface{}          `json:"events,omitempty"`
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

// AccountTransactionsResponse represents the response from the account transactions endpoint
type AccountTransactionsResponse struct {
	Data  []AccountTransaction   `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// TaxReportEntry represents a tax report entry
type TaxReportEntry struct {
	AbsAmount       float64 `json:"abs_amount"`
	Address         string  `json:"address"`
	Amount          float64 `json:"amount"`
	BlockHeight     uint64  `json:"block_height"`
	Direction       string  `json:"direction"`
	Fee             float64 `json:"fee"`
	Otherside       string  `json:"otherside"`
	Time            string  `json:"time"`
	Token           string  `json:"token"`
	TransactionHash string  `json:"transaction_hash"`
	Type            string  `json:"type"`
}

// TaxReportResponse represents the response from the tax report endpoint
type TaxReportResponse struct {
	Data  []TaxReportEntry       `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// AccountsRequestBuilder builds a request to get accounts list
type AccountsRequestBuilder struct {
	service *Service
	height  *uint64
	limit   *int
	offset  *int
	sortBy  *string
}

// GetAccounts creates a new accounts list request builder
func (s *Service) GetAccounts() *AccountsRequestBuilder {
	return &AccountsRequestBuilder{service: s}
}

// Height sets the block height cursor for pagination (optional)
func (b *AccountsRequestBuilder) Height(height uint64) *AccountsRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountsRequestBuilder) Limit(limit int) *AccountsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountsRequestBuilder) Offset(offset int) *AccountsRequestBuilder {
	b.offset = &offset
	return b
}

// SortBy sets the sort field (optional, e.g., "flow_balance")
func (b *AccountsRequestBuilder) SortBy(sortBy string) *AccountsRequestBuilder {
	b.sortBy = &sortBy
	return b
}

// Do executes the accounts list request
func (b *AccountsRequestBuilder) Do(ctx context.Context) (*AccountsResponse, error) {
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
	if b.sortBy != nil {
		query.Set("sort_by", *b.sortBy)
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/account", query)
	if err != nil {
		return nil, err
	}

	var accountsResp AccountsResponse
	if err := b.service.client.DecodeResponse(resp, &accountsResp); err != nil {
		return nil, err
	}

	return &accountsResp, nil
}

// AccountRequestBuilder builds a request to get account details
type AccountRequestBuilder struct {
	service *Service
	address string
}

// GetAccount creates a new account details request builder
func (s *Service) GetAccount() *AccountRequestBuilder {
	return &AccountRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountRequestBuilder) Address(address string) *AccountRequestBuilder {
	b.address = address
	return b
}

// Do executes the account details request
func (b *AccountRequestBuilder) Do(ctx context.Context) (*AccountDetailsResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}

	path := fmt.Sprintf("/flow/v1/account/%s", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var accountResp AccountDetailsResponse
	if err := b.service.client.DecodeResponse(resp, &accountResp); err != nil {
		return nil, err
	}

	return &accountResp, nil
}

// AccountFTsRequestBuilder builds a request to get account FT collections
type AccountFTsRequestBuilder struct {
	service *Service
	address string
	limit   *int
	offset  *int
}

// GetAccountFTs creates a new account FT collections request builder
func (s *Service) GetAccountFTs() *AccountFTsRequestBuilder {
	return &AccountFTsRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountFTsRequestBuilder) Address(address string) *AccountFTsRequestBuilder {
	b.address = address
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountFTsRequestBuilder) Limit(limit int) *AccountFTsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountFTsRequestBuilder) Offset(offset int) *AccountFTsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account FT collections request
func (b *AccountFTsRequestBuilder) Do(ctx context.Context) (*AccountFTCollectionsResponse, error) {
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

	path := fmt.Sprintf("/flow/v1/account/%s/ft", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var collectionsResp AccountFTCollectionsResponse
	if err := b.service.client.DecodeResponse(resp, &collectionsResp); err != nil {
		return nil, err
	}

	return &collectionsResp, nil
}

// AccountFTHoldingsRequestBuilder builds a request to get account FT holdings with statistics
type AccountFTHoldingsRequestBuilder struct {
	service *Service
	address string
	limit   *int
	offset  *int
}

// GetAccountFTHoldings creates a new account FT holdings request builder
func (s *Service) GetAccountFTHoldings() *AccountFTHoldingsRequestBuilder {
	return &AccountFTHoldingsRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountFTHoldingsRequestBuilder) Address(address string) *AccountFTHoldingsRequestBuilder {
	b.address = address
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountFTHoldingsRequestBuilder) Limit(limit int) *AccountFTHoldingsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountFTHoldingsRequestBuilder) Offset(offset int) *AccountFTHoldingsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account FT holdings request
func (b *AccountFTHoldingsRequestBuilder) Do(ctx context.Context) (*FTHoldingResponse, error) {
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

	path := fmt.Sprintf("/flow/v1/account/%s/ft/holding", b.address)
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

// AccountFTTransfersRequestBuilder builds a request to get account FT transfers
type AccountFTTransfersRequestBuilder struct {
	service *Service
	address string
	height  *uint64
	limit   *int
	offset  *int
}

// GetAccountFTTransfers creates a new account FT transfers request builder
func (s *Service) GetAccountFTTransfers() *AccountFTTransfersRequestBuilder {
	return &AccountFTTransfersRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountFTTransfersRequestBuilder) Address(address string) *AccountFTTransfersRequestBuilder {
	b.address = address
	return b
}

// Height sets the block height filter (optional)
func (b *AccountFTTransfersRequestBuilder) Height(height uint64) *AccountFTTransfersRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountFTTransfersRequestBuilder) Limit(limit int) *AccountFTTransfersRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountFTTransfersRequestBuilder) Offset(offset int) *AccountFTTransfersRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account FT transfers request
func (b *AccountFTTransfersRequestBuilder) Do(ctx context.Context) (*TransfersResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}

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

	path := fmt.Sprintf("/flow/v1/account/%s/ft/transfer", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var transfersResp TransfersResponse
	if err := b.service.client.DecodeResponse(resp, &transfersResp); err != nil {
		return nil, err
	}

	return &transfersResp, nil
}

// AccountFTTokenRequestBuilder builds a request to get account's specific FT token
type AccountFTTokenRequestBuilder struct {
	service *Service
	address string
	token   string
	limit   *int
	offset  *int
}

// GetAccountFTToken creates a new account FT token request builder
func (s *Service) GetAccountFTToken() *AccountFTTokenRequestBuilder {
	return &AccountFTTokenRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountFTTokenRequestBuilder) Address(address string) *AccountFTTokenRequestBuilder {
	b.address = address
	return b
}

// Token sets the token identifier (required)
func (b *AccountFTTokenRequestBuilder) Token(token string) *AccountFTTokenRequestBuilder {
	b.token = token
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountFTTokenRequestBuilder) Limit(limit int) *AccountFTTokenRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountFTTokenRequestBuilder) Offset(offset int) *AccountFTTokenRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account FT token request
func (b *AccountFTTokenRequestBuilder) Do(ctx context.Context) (*AccountFungibleTokenResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}
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

	path := fmt.Sprintf("/flow/v1/account/%s/ft/%s", b.address, b.token)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var tokenResp AccountFungibleTokenResponse
	if err := b.service.client.DecodeResponse(resp, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
}

// AccountFTTokenTransfersRequestBuilder builds a request to get account's specific token transfers
type AccountFTTokenTransfersRequestBuilder struct {
	service *Service
	address string
	token   string
	height  *uint64
	limit   *int
	offset  *int
}

// GetAccountFTTokenTransfers creates a new account FT token transfers request builder
func (s *Service) GetAccountFTTokenTransfers() *AccountFTTokenTransfersRequestBuilder {
	return &AccountFTTokenTransfersRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountFTTokenTransfersRequestBuilder) Address(address string) *AccountFTTokenTransfersRequestBuilder {
	b.address = address
	return b
}

// Token sets the token identifier (required)
func (b *AccountFTTokenTransfersRequestBuilder) Token(token string) *AccountFTTokenTransfersRequestBuilder {
	b.token = token
	return b
}

// Height sets the block height filter (optional)
func (b *AccountFTTokenTransfersRequestBuilder) Height(height uint64) *AccountFTTokenTransfersRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountFTTokenTransfersRequestBuilder) Limit(limit int) *AccountFTTokenTransfersRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountFTTokenTransfersRequestBuilder) Offset(offset int) *AccountFTTokenTransfersRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account FT token transfers request
func (b *AccountFTTokenTransfersRequestBuilder) Do(ctx context.Context) (*TransfersResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}
	if b.token == "" {
		return nil, fmt.Errorf("token identifier is required")
	}

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

	path := fmt.Sprintf("/flow/v1/account/%s/ft/%s/transfer", b.address, b.token)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var transfersResp TransfersResponse
	if err := b.service.client.DecodeResponse(resp, &transfersResp); err != nil {
		return nil, err
	}

	return &transfersResp, nil
}

// AccountTaxReportRequestBuilder builds a request to get account tax report
type AccountTaxReportRequestBuilder struct {
	service *Service
	address string
	height  *uint64
	limit   *int
	offset  *int
}

// GetAccountTaxReport creates a new account tax report request builder
func (s *Service) GetAccountTaxReport() *AccountTaxReportRequestBuilder {
	return &AccountTaxReportRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountTaxReportRequestBuilder) Address(address string) *AccountTaxReportRequestBuilder {
	b.address = address
	return b
}

// Height sets the block height filter (optional)
func (b *AccountTaxReportRequestBuilder) Height(height uint64) *AccountTaxReportRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default returns all, max 100)
func (b *AccountTaxReportRequestBuilder) Limit(limit int) *AccountTaxReportRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountTaxReportRequestBuilder) Offset(offset int) *AccountTaxReportRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account tax report request
func (b *AccountTaxReportRequestBuilder) Do(ctx context.Context) (*TaxReportResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}

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

	path := fmt.Sprintf("/flow/v1/account/%s/tax-report", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var taxResp TaxReportResponse
	if err := b.service.client.DecodeResponse(resp, &taxResp); err != nil {
		return nil, err
	}

	return &taxResp, nil
}

// AccountTransactionsRequestBuilder builds a request to get account transactions
type AccountTransactionsRequestBuilder struct {
	service       *Service
	address       string
	height        *uint64
	limit         *int
	offset        *int
	includeEvents *bool
	active        *bool
	from          *string
	to            *string
}

// GetAccountTransactions creates a new account transactions request builder
func (s *Service) GetAccountTransactions() *AccountTransactionsRequestBuilder {
	return &AccountTransactionsRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountTransactionsRequestBuilder) Address(address string) *AccountTransactionsRequestBuilder {
	b.address = address
	return b
}

// Height sets the block height filter (optional)
func (b *AccountTransactionsRequestBuilder) Height(height uint64) *AccountTransactionsRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountTransactionsRequestBuilder) Limit(limit int) *AccountTransactionsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountTransactionsRequestBuilder) Offset(offset int) *AccountTransactionsRequestBuilder {
	b.offset = &offset
	return b
}

// IncludeEvents sets whether to include events in the response (optional, default false)
func (b *AccountTransactionsRequestBuilder) IncludeEvents(include bool) *AccountTransactionsRequestBuilder {
	b.includeEvents = &include
	return b
}

// Active sets whether to get only active transactions (optional)
func (b *AccountTransactionsRequestBuilder) Active(active bool) *AccountTransactionsRequestBuilder {
	b.active = &active
	return b
}

// From sets the start time filter (optional)
func (b *AccountTransactionsRequestBuilder) From(from string) *AccountTransactionsRequestBuilder {
	b.from = &from
	return b
}

// To sets the end time filter (optional)
func (b *AccountTransactionsRequestBuilder) To(to string) *AccountTransactionsRequestBuilder {
	b.to = &to
	return b
}

// Do executes the account transactions request
func (b *AccountTransactionsRequestBuilder) Do(ctx context.Context) (*AccountTransactionsResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}

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
	if b.includeEvents != nil {
		query.Set("include_events", strconv.FormatBool(*b.includeEvents))
	}
	if b.active != nil {
		query.Set("active", strconv.FormatBool(*b.active))
	}
	if b.from != nil {
		query.Set("from", *b.from)
	}
	if b.to != nil {
		query.Set("to", *b.to)
	}

	path := fmt.Sprintf("/flow/v1/account/%s/transaction", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var txResp AccountTransactionsResponse
	if err := b.service.client.DecodeResponse(resp, &txResp); err != nil {
		return nil, err
	}

	return &txResp, nil
}
