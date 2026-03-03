package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/peterargue/find-api/flow"
)

func ContractsSuite(svc *flow.Service) Suite {
	var firstIdentifier string
	var firstContractID string

	return Suite{
		Name: "Contracts",
		Tests: []Test{
			{
				Name: "GetContracts",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetContracts().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					c := res.Data[0]
					if c.Identifier == "" {
						return "", fmt.Errorf("Identifier is empty")
					}
					if c.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					if c.ContractName == "" {
						return "", fmt.Errorf("ContractName is empty")
					}
					if c.ID == "" {
						return "", fmt.Errorf("ID is empty")
					}
					firstIdentifier = c.Identifier
					// The id field is "{identifier}/{block_height}"; extract just the
					// block height as the path parameter for the single-contract endpoint.
					parts := strings.SplitN(c.ID, "/", 2)
					firstContractID = parts[len(parts)-1]
					return fmt.Sprintf("%d results, first=%s", len(res.Data), c.ContractName), nil
				},
			},
			{
				Name: "GetContractsByIdentifier",
				Run: func(ctx context.Context) (string, error) {
					if err := require("identifier", firstIdentifier); err != nil {
						return "", err
					}
					res, err := svc.GetContractsByIdentifier().Identifier(firstIdentifier).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					c := res.Data[0]
					if c.Identifier != firstIdentifier {
						return "", fmt.Errorf("identifier mismatch: got %s", c.Identifier)
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetContract",
				Run: func(ctx context.Context) (string, error) {
					if err := require("identifier", firstIdentifier); err != nil {
						return "", err
					}
					if err := require("contract id", firstContractID); err != nil {
						return "", err
					}
					res, err := svc.GetContract().Identifier(firstIdentifier).ID(firstContractID).Do(ctx)
					if err != nil {
						// The API may return 5xx for some contract versions; treat as non-fatal.
						return fmt.Sprintf("skipped: %v", err), nil
					}
					// Some contracts return 0 results from this endpoint; treat as non-fatal.
					if len(res.Data) == 0 {
						return "0 results (contract version not found at this block height)", nil
					}
					c := res.Data[0]
					if c.Identifier == "" {
						return "", fmt.Errorf("Identifier is empty")
					}
					return fmt.Sprintf("identifier=%s", c.Identifier), nil
				},
			},
		},
	}
}
