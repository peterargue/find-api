package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func TransactionSuite(svc *flow.Service) Suite {
	var firstTxID string

	return Suite{
		Name: "Transactions",
		Tests: []Test{
			{
				Name: "GetTransactions",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetTransactions().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					tx := res.Data[0]
					if tx.ID == "" {
						return "", fmt.Errorf("ID is empty")
					}
					if tx.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					if tx.Status == "" {
						return "", fmt.Errorf("Status is empty")
					}
					firstTxID = tx.ID
					return fmt.Sprintf("%d results, first=%s", len(res.Data), tx.ID), nil
				},
			},
			{
				Name: "GetTransaction",
				Run: func(ctx context.Context) (string, error) {
					if err := require("tx id", firstTxID); err != nil {
						return "", err
					}
					res, err := svc.GetTransaction().ID(firstTxID).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					tx := res.Data[0]
					if tx.ID == "" {
						return "", fmt.Errorf("ID is empty")
					}
					if tx.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					if tx.Status == "" {
						return "", fmt.Errorf("Status is empty")
					}
					return fmt.Sprintf("id=%s, status=%s", tx.ID, tx.Status), nil
				},
			},
			{
				Name: "GetScheduledTransactions",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetScheduledTransactions().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					// Scheduled transactions may be empty if none exist.
					if len(res.Data) == 0 {
						return "0 results", nil
					}
					tx := res.Data[0]
					if tx.ID == "" {
						return "", fmt.Errorf("ID is empty")
					}
					if tx.Status == "" {
						return "", fmt.Errorf("Status is empty")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
		},
	}
}
