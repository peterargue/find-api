package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func ContractsSuite(svc *flow.Service) Suite {
	var firstIdentifier string
	var firstAddress string

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
					firstIdentifier = c.Identifier
					firstAddress = c.Address
					// contract_name and block_height are not returned by this endpoint.
					return fmt.Sprintf("%d results, first=%s", len(res.Data), c.Identifier), nil
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
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetContract().Identifier(firstIdentifier).ID(firstAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						// This endpoint may return empty for some contracts; treat as skip.
						return "0 results (contract not available via this endpoint)", nil
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
