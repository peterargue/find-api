package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func NFTSuite(svc *flow.Service) Suite {
	var firstNFTType string
	var firstNFTId string
	var firstNFTAddress string

	return Suite{
		Name: "NFT",
		Tests: []Test{
			{
				Name: "GetNFTCollections",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetNFTCollections().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					c := res.Data[0]
					if c.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					if c.Name == "" {
						return "", fmt.Errorf("Name is empty")
					}
					firstNFTType = c.NFTType
					return fmt.Sprintf("%d results, first=%s", len(res.Data), c.Name), nil
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
					c := res.Data[0]
					if c.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					return fmt.Sprintf("items=%d, holders=%d", c.ItemCount, c.HolderCount), nil
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
					if t.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					if t.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					if t.TransactionID == "" {
						return "", fmt.Errorf("TransactionID is empty")
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
					if h.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					if h.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					return fmt.Sprintf("%d results, top=%s", len(res.Data), h.Address), nil
				},
			},
			{
				Name: "GetNFTItem",
				Run: func(ctx context.Context) (string, error) {
					if err := require("nft type", firstNFTType); err != nil {
						return "", err
					}
					if err := require("nft id", firstNFTId); err != nil {
						return "", err
					}
					res, err := svc.GetNFTItem().NFTType(firstNFTType).ID(firstNFTId).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					n := res.Data[0]
					if n.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					if n.NFTId == "" {
						return "", fmt.Errorf("NFTId is empty")
					}
					return fmt.Sprintf("id=%s", n.NFTId), nil
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
					c := res.Data[0]
					if c.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
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
					n := res.Data[0]
					if n.NFTType == "" {
						return "", fmt.Errorf("NFTType is empty")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
		},
	}
}
