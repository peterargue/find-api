package evm

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var tokenCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "token <address>",
		Short: "Get an EVM token by contract address",
		Args:  cobra.ExactArgs(1),
	},
	Run: runToken,
}

type evmTokenResult struct {
	token flow.EvmToken
}

func (r *evmTokenResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "name:         %s\n", r.token.Name)
	fmt.Fprintf(&buf, "symbol:       %s\n", r.token.Symbol)
	fmt.Fprintf(&buf, "type:         %s\n", r.token.Type)
	fmt.Fprintf(&buf, "address:      %s\n", r.token.ContractAddressHash)
	fmt.Fprintf(&buf, "decimals:     %d\n", r.token.Decimals)
	fmt.Fprintf(&buf, "holders:      %d\n", r.token.Holders)
	fmt.Fprintf(&buf, "transfers:    %d\n", r.token.Transfers)
	fmt.Fprintf(&buf, "total_supply: %s\n", r.token.TotalSupply)
	return buf.String()
}

func (r *evmTokenResult) Oneliner() string {
	return fmt.Sprintf("%s (%s)", r.token.Name, r.token.Symbol)
}

func (r *evmTokenResult) JSON() any { return r.token }

func runToken(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetEvmToken().Address(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("EVM token not found")
	}
	return &evmTokenResult{token: resp.Data[0]}, nil
}
