package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func EvmSuite(svc *flow.Service) Suite {
	var firstEvmTokenAddress string
	var firstEvmTxHash string

	return Suite{
		Name: "EVM",
		Tests: []Test{
			{
				Name: "GetEvmTokens",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetEvmTokens().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					t := res.Data[0]
					if t.ContractAddressHash == "" {
						return "", fmt.Errorf("ContractAddressHash is empty")
					}
					if t.Symbol == "" {
						return "", fmt.Errorf("Symbol is empty")
					}
					if t.Type == "" {
						return "", fmt.Errorf("Type is empty")
					}
					firstEvmTokenAddress = t.ContractAddressHash
					return fmt.Sprintf("%d results, first=%s", len(res.Data), t.Symbol), nil
				},
			},
			{
				Name: "GetEvmToken",
				Run: func(ctx context.Context) (string, error) {
					if err := require("evm token address", firstEvmTokenAddress); err != nil {
						return "", err
					}
					res, err := svc.GetEvmToken().Address(firstEvmTokenAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					t := res.Data[0]
					if t.Symbol == "" {
						return "", fmt.Errorf("Symbol is empty")
					}
					return fmt.Sprintf("symbol=%s, holders=%d", t.Symbol, t.Holders), nil
				},
			},
			{
				Name: "GetEvmTransactions",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetEvmTransactions().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					tx := res.Data[0]
					if tx.Hash == "" {
						return "", fmt.Errorf("Hash is empty")
					}
					if tx.BlockNumber == 0 {
						return "", fmt.Errorf("BlockNumber is zero")
					}
					firstEvmTxHash = tx.Hash
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetEvmTransaction",
				Run: func(ctx context.Context) (string, error) {
					if err := require("evm tx hash", firstEvmTxHash); err != nil {
						return "", err
					}
					res, err := svc.GetEvmTransaction().Hash(firstEvmTxHash).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						dumpJSON("EvmTransactionResponse", res)
						return "", fmt.Errorf("empty response")
					}
					tx := res.Data[0]
					if tx.Hash == "" {
						dumpJSON("EvmTransaction[0]", tx)
						return "", fmt.Errorf("Hash is empty")
					}
					if tx.BlockNumber == 0 {
						return "", fmt.Errorf("BlockNumber is zero")
					}
					if tx.Status == "" {
						return "", fmt.Errorf("Status is empty")
					}
					return fmt.Sprintf("hash=%s, status=%s", tx.Hash, tx.Status), nil
				},
			},
		},
	}
}
