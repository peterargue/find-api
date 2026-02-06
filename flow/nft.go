package flow

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

// NFTCollection represents an NFT collection
type NFTCollection struct {
	Address      string `json:"address"`
	Banner       string `json:"banner"`
	ContractName string `json:"contract_name"`
	Description  string `json:"description"`
	Display      string `json:"display"`
	External     string `json:"external"`
	FlowtyID     string `json:"flowty_id"`
	Logo         string `json:"logo"`
	Name         string `json:"name"`
	NFTType      string `json:"nft_type"`
	Path         string `json:"path"`
	Twitter      string `json:"twitter"`
	Website      string `json:"website"`
}

// NFTCollectionResponse represents the response from the NFT collections list endpoint
type NFTCollectionResponse struct {
	Data  []NFTCollection        `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NFTCollectionDetails represents detailed NFT collection information
type NFTCollectionDetails struct {
	NFTCollection
	HolderCount int `json:"holder_count"`
	ItemCount   int `json:"item_count"`
}

// NFTCollectionDetailsResponse represents the response from the NFT collection details endpoint
type NFTCollectionDetailsResponse struct {
	Data  []NFTCollectionDetails `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NFTTransfer represents an NFT transfer
type NFTTransfer struct {
	Address         string `json:"address"`
	BlockHeight     uint64 `json:"block_height"`
	Direction       string `json:"direction"`
	NFTId           string `json:"nft_id"`
	NFTType         string `json:"nft_type"`
	Receiver        string `json:"receiver"`
	Sender          string `json:"sender"`
	Timestamp       string `json:"timestamp"`
	TransactionHash string `json:"transaction_hash"`
	TransactionID   string `json:"transaction_id"`
}

// NFTTransfersResponse represents the response from the NFT transfers endpoint
type NFTTransfersResponse struct {
	Data  []NFTTransfer          `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NFTHolding represents an NFT holding
type NFTHolding struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
	NFTType string `json:"nft_type"`
}

// NFTHoldingResponse represents the response from the NFT holdings endpoint
type NFTHoldingResponse struct {
	Data  []NFTHolding           `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NFT represents detailed NFT information
type NFT struct {
	Address     string                 `json:"address"`
	BlockHeight uint64                 `json:"block_height"`
	ID          string                 `json:"id"`
	Metadata    map[string]interface{} `json:"metadata"`
	Name        string                 `json:"name"`
	NFTId       string                 `json:"nft_id"`
	NFTType     string                 `json:"nft_type"`
	Owner       string                 `json:"owner"`
	Thumbnail   string                 `json:"thumbnail"`
}

// NFTDetailsResponse represents the response from the NFT details endpoint
type NFTDetailsResponse struct {
	Data  []NFT                  `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// AccountNFTCollection represents an NFT collection summary for an account
type AccountNFTCollection struct {
	Banner   string `json:"banner"`
	Logo     string `json:"logo"`
	Name     string `json:"name"`
	NFTCount int    `json:"nft_count"`
	NFTType  string `json:"nft_type"`
	Owner    string `json:"owner"`
}

// AccountNFTCollectionsResponse represents the response from account NFT collections endpoint
type AccountNFTCollectionsResponse struct {
	Data  []AccountNFTCollection `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// AccountNFT represents an NFT owned by an account
type AccountNFT struct {
	Address     string                 `json:"address"`
	BlockHeight uint64                 `json:"block_height"`
	ID          string                 `json:"id"`
	Metadata    map[string]interface{} `json:"metadata"`
	Name        string                 `json:"name"`
	NFTId       string                 `json:"nft_id"`
	NFTType     string                 `json:"nft_type"`
	Owner       string                 `json:"owner"`
	Thumbnail   string                 `json:"thumbnail"`
	Valid       bool                   `json:"valid"`
}

// AccountNFTResponse represents the response from account NFT endpoint
type AccountNFTResponse struct {
	Data  []AccountNFT           `json:"data"`
	Links map[string]string      `json:"_links"`
	Meta  map[string]interface{} `json:"_meta"`
	Error interface{}            `json:"error,omitempty"`
}

// NFTCollectionsRequestBuilder builds a request to get NFT collections
type NFTCollectionsRequestBuilder struct {
	service *Service
	limit   *int
	offset  *int
}

// GetNFTCollections creates a new NFT collections request builder
func (s *Service) GetNFTCollections() *NFTCollectionsRequestBuilder {
	return &NFTCollectionsRequestBuilder{service: s}
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *NFTCollectionsRequestBuilder) Limit(limit int) *NFTCollectionsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *NFTCollectionsRequestBuilder) Offset(offset int) *NFTCollectionsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the NFT collections request
func (b *NFTCollectionsRequestBuilder) Do(ctx context.Context) (*NFTCollectionResponse, error) {
	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/nft", query)
	if err != nil {
		return nil, err
	}

	var nftResp NFTCollectionResponse
	if err := b.service.client.DecodeResponse(resp, &nftResp); err != nil {
		return nil, err
	}

	return &nftResp, nil
}

// NFTCollectionRequestBuilder builds a request to get NFT collection details
type NFTCollectionRequestBuilder struct {
	service *Service
	nftType string
}

// GetNFTCollection creates a new NFT collection details request builder
func (s *Service) GetNFTCollection() *NFTCollectionRequestBuilder {
	return &NFTCollectionRequestBuilder{service: s}
}

// NFTType sets the NFT collection type (required)
func (b *NFTCollectionRequestBuilder) NFTType(nftType string) *NFTCollectionRequestBuilder {
	b.nftType = nftType
	return b
}

// Do executes the NFT collection details request
func (b *NFTCollectionRequestBuilder) Do(ctx context.Context) (*NFTCollectionDetailsResponse, error) {
	if b.nftType == "" {
		return nil, fmt.Errorf("NFT type is required")
	}

	path := fmt.Sprintf("/flow/v1/nft/%s", b.nftType)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var nftResp NFTCollectionDetailsResponse
	if err := b.service.client.DecodeResponse(resp, &nftResp); err != nil {
		return nil, err
	}

	return &nftResp, nil
}

// NFTTransfersRequestBuilder builds a request to get NFT transfers
type NFTTransfersRequestBuilder struct {
	service *Service
	address *string
	height  *uint64
	limit   *int
	nftID   *int
	nftType *string
	offset  *int
}

// GetNFTTransfers creates a new NFT transfers request builder
func (s *Service) GetNFTTransfers() *NFTTransfersRequestBuilder {
	return &NFTTransfersRequestBuilder{service: s}
}

// Address sets the address filter (optional)
func (b *NFTTransfersRequestBuilder) Address(address string) *NFTTransfersRequestBuilder {
	b.address = &address
	return b
}

// Height sets the block height filter (optional)
func (b *NFTTransfersRequestBuilder) Height(height uint64) *NFTTransfersRequestBuilder {
	b.height = &height
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *NFTTransfersRequestBuilder) Limit(limit int) *NFTTransfersRequestBuilder {
	b.limit = &limit
	return b
}

// NFTId sets the NFT ID filter (optional)
func (b *NFTTransfersRequestBuilder) NFTId(nftID int) *NFTTransfersRequestBuilder {
	b.nftID = &nftID
	return b
}

// NFTType sets the NFT type filter (optional)
func (b *NFTTransfersRequestBuilder) NFTType(nftType string) *NFTTransfersRequestBuilder {
	b.nftType = &nftType
	return b
}

// Offset sets the pagination offset (optional)
func (b *NFTTransfersRequestBuilder) Offset(offset int) *NFTTransfersRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the NFT transfers request
func (b *NFTTransfersRequestBuilder) Do(ctx context.Context) (*NFTTransfersResponse, error) {
	query := url.Values{}
	if b.address != nil {
		query.Set("address", *b.address)
	}
	if b.height != nil {
		query.Set("height", strconv.FormatUint(*b.height, 10))
	}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.nftID != nil {
		query.Set("nft_id", strconv.Itoa(*b.nftID))
	}
	if b.nftType != nil {
		query.Set("nft_type", *b.nftType)
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, "/flow/v1/nft/transfer", query)
	if err != nil {
		return nil, err
	}

	var transfersResp NFTTransfersResponse
	if err := b.service.client.DecodeResponse(resp, &transfersResp); err != nil {
		return nil, err
	}

	return &transfersResp, nil
}

// NFTHoldingsRequestBuilder builds a request to get NFT holdings
type NFTHoldingsRequestBuilder struct {
	service *Service
	nftType string
	limit   *int
	offset  *int
}

// GetNFTHoldings creates a new NFT holdings request builder
func (s *Service) GetNFTHoldings() *NFTHoldingsRequestBuilder {
	return &NFTHoldingsRequestBuilder{service: s}
}

// NFTType sets the NFT type (required)
func (b *NFTHoldingsRequestBuilder) NFTType(nftType string) *NFTHoldingsRequestBuilder {
	b.nftType = nftType
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *NFTHoldingsRequestBuilder) Limit(limit int) *NFTHoldingsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *NFTHoldingsRequestBuilder) Offset(offset int) *NFTHoldingsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the NFT holdings request
func (b *NFTHoldingsRequestBuilder) Do(ctx context.Context) (*NFTHoldingResponse, error) {
	if b.nftType == "" {
		return nil, fmt.Errorf("NFT type is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}

	path := fmt.Sprintf("/flow/v1/nft/%s/holding", b.nftType)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var holdingsResp NFTHoldingResponse
	if err := b.service.client.DecodeResponse(resp, &holdingsResp); err != nil {
		return nil, err
	}

	return &holdingsResp, nil
}

// NFTItemRequestBuilder builds a request to get NFT item details
type NFTItemRequestBuilder struct {
	service *Service
	nftType string
	id      string
}

// GetNFTItem creates a new NFT item details request builder
func (s *Service) GetNFTItem() *NFTItemRequestBuilder {
	return &NFTItemRequestBuilder{service: s}
}

// NFTType sets the NFT type (required)
func (b *NFTItemRequestBuilder) NFTType(nftType string) *NFTItemRequestBuilder {
	b.nftType = nftType
	return b
}

// ID sets the NFT ID (required)
func (b *NFTItemRequestBuilder) ID(id string) *NFTItemRequestBuilder {
	b.id = id
	return b
}

// Do executes the NFT item details request
func (b *NFTItemRequestBuilder) Do(ctx context.Context) (*NFTDetailsResponse, error) {
	if b.nftType == "" {
		return nil, fmt.Errorf("NFT type is required")
	}
	if b.id == "" {
		return nil, fmt.Errorf("NFT ID is required")
	}

	path := fmt.Sprintf("/flow/v1/nft/%s/item/%s", b.nftType, b.id)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var nftResp NFTDetailsResponse
	if err := b.service.client.DecodeResponse(resp, &nftResp); err != nil {
		return nil, err
	}

	return &nftResp, nil
}

// AccountNFTCollectionsRequestBuilder builds a request to get account NFT collections
type AccountNFTCollectionsRequestBuilder struct {
	service *Service
	address string
	limit   *int
	offset  *int
}

// GetAccountNFTCollections creates a new account NFT collections request builder
func (s *Service) GetAccountNFTCollections() *AccountNFTCollectionsRequestBuilder {
	return &AccountNFTCollectionsRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountNFTCollectionsRequestBuilder) Address(address string) *AccountNFTCollectionsRequestBuilder {
	b.address = address
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountNFTCollectionsRequestBuilder) Limit(limit int) *AccountNFTCollectionsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountNFTCollectionsRequestBuilder) Offset(offset int) *AccountNFTCollectionsRequestBuilder {
	b.offset = &offset
	return b
}

// Do executes the account NFT collections request
func (b *AccountNFTCollectionsRequestBuilder) Do(ctx context.Context) (*AccountNFTCollectionsResponse, error) {
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

	path := fmt.Sprintf("/flow/v1/account/%s/nft", b.address)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var collectionsResp AccountNFTCollectionsResponse
	if err := b.service.client.DecodeResponse(resp, &collectionsResp); err != nil {
		return nil, err
	}

	return &collectionsResp, nil
}

// AccountNFTsRequestBuilder builds a request to get account NFTs by collection
type AccountNFTsRequestBuilder struct {
	service   *Service
	address   string
	nftType   string
	limit     *int
	offset    *int
	validOnly *bool
	sortBy    *string
}

// GetAccountNFTs creates a new account NFTs request builder
func (s *Service) GetAccountNFTs() *AccountNFTsRequestBuilder {
	return &AccountNFTsRequestBuilder{service: s}
}

// Address sets the account address (required)
func (b *AccountNFTsRequestBuilder) Address(address string) *AccountNFTsRequestBuilder {
	b.address = address
	return b
}

// NFTType sets the NFT collection type (required)
func (b *AccountNFTsRequestBuilder) NFTType(nftType string) *AccountNFTsRequestBuilder {
	b.nftType = nftType
	return b
}

// Limit sets the number of records to return (optional, default 25, max 100)
func (b *AccountNFTsRequestBuilder) Limit(limit int) *AccountNFTsRequestBuilder {
	b.limit = &limit
	return b
}

// Offset sets the pagination offset (optional)
func (b *AccountNFTsRequestBuilder) Offset(offset int) *AccountNFTsRequestBuilder {
	b.offset = &offset
	return b
}

// ValidOnly sets whether to return only valid NFTs (optional, default false)
func (b *AccountNFTsRequestBuilder) ValidOnly(validOnly bool) *AccountNFTsRequestBuilder {
	b.validOnly = &validOnly
	return b
}

// SortBy sets the sort order (optional, e.g., "asc" or "desc")
func (b *AccountNFTsRequestBuilder) SortBy(sortBy string) *AccountNFTsRequestBuilder {
	b.sortBy = &sortBy
	return b
}

// Do executes the account NFTs request
func (b *AccountNFTsRequestBuilder) Do(ctx context.Context) (*AccountNFTResponse, error) {
	if b.address == "" {
		return nil, fmt.Errorf("account address is required")
	}
	if b.nftType == "" {
		return nil, fmt.Errorf("NFT type is required")
	}

	query := url.Values{}
	if b.limit != nil {
		query.Set("limit", strconv.Itoa(*b.limit))
	}
	if b.offset != nil {
		query.Set("offset", strconv.Itoa(*b.offset))
	}
	if b.validOnly != nil {
		query.Set("valid_only", strconv.FormatBool(*b.validOnly))
	}
	if b.sortBy != nil {
		query.Set("sort_by", *b.sortBy)
	}

	path := fmt.Sprintf("/flow/v1/account/%s/nft/%s", b.address, b.nftType)
	resp, err := b.service.client.DoRequest(ctx, http.MethodGet, path, query)
	if err != nil {
		return nil, err
	}

	var nftResp AccountNFTResponse
	if err := b.service.client.DecodeResponse(resp, &nftResp); err != nil {
		return nil, err
	}

	return &nftResp, nil
}
