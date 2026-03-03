package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func BlocksSuite(svc *flow.Service) Suite {
	var firstHeight uint64

	return Suite{
		Name: "Blocks",
		Tests: []Test{
			{
				Name: "GetBlocks",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetBlocks().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					b := res.Data[0]
					if b.Height == 0 {
						return "", fmt.Errorf("Height is zero")
					}
					if b.ID == "" {
						return "", fmt.Errorf("ID is empty")
					}
					if b.Timestamp == "" {
						return "", fmt.Errorf("Timestamp is empty")
					}
					firstHeight = b.Height
					return fmt.Sprintf("%d results, height=%d", len(res.Data), b.Height), nil
				},
			},
			{
				Name: "GetBlock",
				Run: func(ctx context.Context) (string, error) {
					if err := requireUint64("height", firstHeight); err != nil {
						return "", err
					}
					res, err := svc.GetBlock().Height(firstHeight).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					b := res.Data[0]
					if b.Height != firstHeight {
						return "", fmt.Errorf("unexpected height %d, want %d", b.Height, firstHeight)
					}
					return fmt.Sprintf("height=%d, txs=%d", b.Height, b.Tx), nil
				},
			},
			{
				Name: "GetBlockServiceEvents",
				Run: func(ctx context.Context) (string, error) {
					if err := requireUint64("height", firstHeight); err != nil {
						return "", err
					}
					res, err := svc.GetBlockServiceEvents().Height(firstHeight).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						// Service events only appear at epoch boundaries; empty is valid.
						return "0 results (not an epoch boundary)", nil
					}
					e := res.Data[0]
					if e.Name == "" {
						return "", fmt.Errorf("Name is empty")
					}
					return fmt.Sprintf("%d results, name=%s", len(res.Data), e.Name), nil
				},
			},
			{
				Name: "GetBlockTransactions",
				Run: func(ctx context.Context) (string, error) {
					if err := requireUint64("height", firstHeight); err != nil {
						return "", err
					}
					res, err := svc.GetBlockTransactions().Height(firstHeight).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "0 results (empty block)", nil
					}
					tx := res.Data[0]
					if tx.TransactionID == "" {
						return "", fmt.Errorf("TransactionID is empty")
					}
					if tx.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
		},
	}
}
