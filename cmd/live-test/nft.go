package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/peterargue/find-api/flow"
)

// nftTypeFromCollection constructs the Cadence NFT type identifier from address
// and contract_name when the nft_type field is empty in the API response.
// Standard form: A.<address_without_0x>.<ContractName>.NFT
func nftTypeFromCollection(c flow.NFTCollection) string {
	if c.NFTType != "" {
		return c.NFTType
	}
	if c.Address == "" || c.ContractName == "" {
		return ""
	}
	addr := strings.TrimPrefix(c.Address, "0x")
	return fmt.Sprintf("A.%s.%s.NFT", addr, c.ContractName)
}

func NFTSuite(svc *flow.Service) Suite {
	var firstNFTType string
	var firstNFTId int64
	var firstNFTAddress string

	return Suite{
		Name: "NFT",
		Tests: []Test{
			{
				Name: "GetNFTCollections",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetNFTCollections().Limit(25).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					// Find the first collection with a usable nft_type AND item count > 0.
					// We probe up to 5 candidates to avoid long startup time.
					probed := 0
					for _, c := range res.Data {
						if c.Name == "" {
							continue
						}
						nftType := nftTypeFromCollection(c)
						if nftType == "" {
							continue
						}
						detail, err := svc.GetNFTCollection().NFTType(nftType).Do(ctx)
						if err != nil || len(detail.Data) == 0 {
							probed++
							if probed >= 5 {
								break
							}
							continue
						}
						if detail.Data[0].ItemCount > 0 {
							firstNFTType = nftType
							return fmt.Sprintf("%d results, using %s (items=%d)", len(res.Data), c.Name, detail.Data[0].ItemCount), nil
						}
						probed++
						if probed >= 5 {
							break
						}
					}
					// Fall back to first usable type even without items.
					for _, c := range res.Data {
						if c.Name == "" {
							continue
						}
						nftType := nftTypeFromCollection(c)
						if nftType != "" {
							firstNFTType = nftType
							return fmt.Sprintf("%d results, first=%s (no active collections found)", len(res.Data), c.Name), nil
						}
					}
					return "", fmt.Errorf("no usable collection found in %d results", len(res.Data))
				},
			},
			{
				Name: "GetNFTCollection",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					res, err := svc.GetNFTCollection().NFTType(firstNFTType).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					return fmt.Sprintf("items=%d, holders=%d", res.Data[0].ItemCount, res.Data[0].HolderCount), nil
				},
			},
			{
				Name: "GetNFTTransfers",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					res, err := svc.GetNFTTransfers().NFTType(firstNFTType).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					t := res.Data[0]
					if t.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					if t.NFTId == 0 {
						return "", fmt.Errorf("NFTId is zero")
					}
					if t.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					firstNFTId = t.NFTId
					firstNFTAddress = t.Address
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetNFTHoldings",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					res, err := svc.GetNFTHoldings().NFTType(firstNFTType).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					h := res.Data[0]
					// NFTType is the only reliably populated field for some collections;
					// Address may be empty when the collection has no active holders.
					if h.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetNFTItem",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					if firstNFTId == 0 {
						return "", fmt.Errorf("prerequisite missing: no nft id from previous test")
					}
					idStr := fmt.Sprintf("%d", firstNFTId)
					res, err := svc.GetNFTItem().NFTType(firstNFTType).ID(idStr).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					n := res.Data[0]
					if n.NFTId == 0 {
						return "", fmt.Errorf("NFTId is zero")
					}
					return fmt.Sprintf("id=%d", n.NFTId), nil
				},
			},
			{
				Name: "GetAccountNFTCollections",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft address", firstNFTAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountNFTCollections().Address(firstNFTAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					return fmt.Sprintf("%d collections, owner=%s", len(res.Data), firstNFTAddress), nil
				},
			},
			{
				Name: "GetAccountNFTs",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft address", firstNFTAddress); err != nil {
						return "", err
					}
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					res, err := svc.GetAccountNFTs().Address(firstNFTAddress).NFTType(firstNFTType).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
		},
	}
}
