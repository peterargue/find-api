package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func FTSuite(svc *flow.Service) Suite {
	var firstToken string
	var firstHolderAddress string

	return Suite{
		Name: "FT",
		Tests: []Test{
			{
				Name: "GetFTs",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetFTs().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					ft := res.Data[0]
					if ft.Token == "" {
						return "", fmt.Errorf("Token is empty")
					}
					if ft.Symbol == "" {
						return "", fmt.Errorf("Symbol is empty")
					}
					if ft.Name == "" {
						return "", fmt.Errorf("Name is empty")
					}
					firstToken = ft.Token
					return fmt.Sprintf("%d results, first=%s", len(res.Data), ft.Symbol), nil
				},
			},
			{
				Name: "GetFT",
				Run: func(ctx context.Context) (string, error) {
					if err := require("token", firstToken); err != nil {
						return "", err
					}
					res, err := svc.GetFT().Token(firstToken).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					ft := res.Data[0]
					if ft.Symbol == "" {
						return "", fmt.Errorf("Symbol is empty")
					}
					return fmt.Sprintf("symbol=%s, owners=%d", ft.Symbol, ft.Stats.OwnerCounts), nil
				},
			},
			{
				Name: "GetFTTransfers",
				Run: func(ctx context.Context) (string, error) {
					if err := require("token", firstToken); err != nil {
						return "", err
					}
					res, err := svc.GetFTTransfers().Token(firstToken).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					t := res.Data[0]
					if t.Amount == 0 {
						return "", fmt.Errorf("Amount is zero")
					}
					if t.Direction == "" {
						return "", fmt.Errorf("Direction is empty")
					}
					if t.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetFTHoldings",
				Run: func(ctx context.Context) (string, error) {
					if err := require("token", firstToken); err != nil {
						return "", err
					}
					res, err := svc.GetFTHoldings().Token(firstToken).Limit(5).Do(ctx)
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
					if h.Balance == 0 {
						return "", fmt.Errorf("Balance is zero")
					}
					firstHolderAddress = h.Address
					return fmt.Sprintf("%d results, top holder=%s", len(res.Data), h.Address), nil
				},
			},
			{
				Name: "GetFTAccountToken",
				Run: func(ctx context.Context) (string, error) {
					if err := require("token", firstToken); err != nil {
						return "", err
					}
					if err := require("holder address", firstHolderAddress); err != nil {
						return "", err
					}
					res, err := svc.GetFTAccountToken().Token(firstToken).Address(firstHolderAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					v := res.Data[0]
					if v.Balance == 0 {
						return "", fmt.Errorf("Balance is zero")
					}
					if v.Token == "" {
						return "", fmt.Errorf("Token is empty")
					}
					return fmt.Sprintf("balance=%g", v.Balance), nil
				},
			},
		},
	}
}
