package main

import (
	"context"
	"fmt"

	"github.com/peterargue/find-api/flow"
)

func NodesSuite(svc *flow.Service) Suite {
	var firstNodeID string

	return Suite{
		Name: "Nodes",
		Tests: []Test{
			{
				Name: "GetNodes",
				Run: func(ctx context.Context) (string, error) {
					res, err := svc.GetNodes().Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					n := res.Data[0]
					if n.NodeID == "" {
						return "", fmt.Errorf("NodeID is empty")
					}
					if n.Role == "" {
						return "", fmt.Errorf("Role is empty")
					}
					firstNodeID = n.NodeID
					return fmt.Sprintf("%d results, first role=%s", len(res.Data), n.Role), nil
				},
			},
			{
				Name: "GetNode",
				Run: func(ctx context.Context) (string, error) {
					if err := require("node id", firstNodeID); err != nil {
						return "", err
					}
					res, err := svc.GetNode().NodeID(firstNodeID).Do(ctx)
					if err != nil {
						return "", err
					}
					if len(res.Data) == 0 {
						return "", fmt.Errorf("empty response")
					}
					n := res.Data[0]
					if n.NodeID == "" {
						return "", fmt.Errorf("NodeID is empty")
					}
					return fmt.Sprintf("id=%s, role=%s", n.NodeID, n.Role), nil
				},
			},
			{
				Name: "GetNodeDelegationRewards",
				Run: func(ctx context.Context) (string, error) {
					if err := require("node id", firstNodeID); err != nil {
						return "", err
					}
					res, err := svc.GetNodeDelegationRewards().NodeID(firstNodeID).Limit(5).Do(ctx)
					if err != nil {
						return "", err
					}
					// Delegation rewards may be empty for nodes with no delegators.
					if len(res.Data) == 0 {
						return "0 results", nil
					}
					r := res.Data[0]
					if r.NodeID == "" {
						return "", fmt.Errorf("NodeID is empty")
					}
					if r.BlockHeight == 0 {
						return "", fmt.Errorf("BlockHeight is zero")
					}
					return fmt.Sprintf("%d results", len(res.Data)), nil
				},
			},
		},
	}
}
