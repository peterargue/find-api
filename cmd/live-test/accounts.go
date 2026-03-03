package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func AccountsSuite(svc *flow.Service) Suite {
	var firstAddress string
	var firstFTToken string

	return Suite{
		Name: "Accounts",
		Tests: []Test{
			{
				Name: "GetAccounts",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetAccounts().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					a := res.Data[0]
					if a.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					if a.Height == 0 {
						return "", fmt.Errorf("Height is zero")
					}
					firstAddress = a.Address
					return fmt.Sprintf("%d results, first=%s", len(res.Data), a.Address), nil
				},
			},
			{
				Name: "GetAccount",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccount().Address(firstAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					a := res.Data[0]
					if a.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					return fmt.Sprintf("address=%s, balance=%g", a.Address, a.FlowBalance), nil
				},
			},
			{
				Name: "GetAccountFTs",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountFTs().Address(firstAddress).Do(ctx)
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
					if ft.Address == "" {
						return "", fmt.Errorf("Address is empty")
					}
					firstFTToken = ft.Token
					return fmt.Sprintf("%d FT collections", len(res.Data)), nil
				},
			},
			{
				Name: "GetAccountFTHoldings",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountFTHoldings().Address(firstAddress).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					h := res.Data[0]
					if h.Token == "" {
						return "", fmt.Errorf("Token is empty")
					}
					return fmt.Sprintf("%d holdings", len(res.Data)), nil
				},
			},
			{
				Name: "GetAccountFTTransfers",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountFTTransfers().Address(firstAddress).Limit(5).Do(ctx)
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
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetAccountFTToken",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					if err := require("FT token", firstFTToken); err != nil {
						return "", err
					}
					res, err := svc.GetAccountFTToken().Address(firstAddress).Token(firstFTToken).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					v := res.Data[0]
					if v.Token == "" {
						return "", fmt.Errorf("Token is empty")
					}
					return fmt.Sprintf("balance=%g", v.Balance), nil
				},
			},
			{
				Name: "GetAccountFTTokenTransfers",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					if err := require("FT token", firstFTToken); err != nil {
						return "", err
					}
					res, err := svc.GetAccountFTTokenTransfers().Address(firstAddress).Token(firstFTToken).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					t := res.Data[0]
					if t.Direction == "" {
						return "", fmt.Errorf("Direction is empty")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
			{
				Name: "GetAccountTaxReport",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountTaxReport().Address(firstAddress).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					// Tax report may legitimately be empty for some accounts.
					return fmt.Sprintf("%d entries", len(res.Data)), nil
				},
			},
			{
				Name: "GetAccountTransactions",
				Run: func(ctx context.Context) (string, error) {
					if err := require("address", firstAddress); err != nil {
						return "", err
					}
					res, err := svc.GetAccountTransactions().Address(firstAddress).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
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
