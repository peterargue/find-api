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
					// Set chain variables before optional field checks so subsequent
					// tests can run even if ContractName is missing.
					firstIdentifier = c.Identifier
					firstAddress = c.Address
					if c.ContractName == "" {
						dumpJSON("Contract[0]", c)
						return "", fmt.Errorf("ContractName is empty")
					}
					if c.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
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
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetContract().Identifier(firstIdentifier).ID(firstAddress).Do(ctx)
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
					return fmt.Sprintf("identifier=%s", c.Identifier), nil
				},
			},
		},
	}
}
