package ft

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "get <token>",
		Short: "Get fungible token details",
		Args:  cobra.ExactArgs(1),
	},
	Run: runGet,
}

type ftGetResult struct {
	token flow.FungibleTokenDetails
}

func (r *ftGetResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Name:          %s\n", r.token.Name)
	fmt.Fprintf(&buf, "Symbol:        %s\n", r.token.Symbol)
	fmt.Fprintf(&buf, "Address:       %s\n", r.token.Address)
	fmt.Fprintf(&buf, "Contract:      %s\n", r.token.ContractName)
	fmt.Fprintf(&buf, "Supply:        %g\n", r.token.CirculatingSupply)
	fmt.Fprintf(&buf, "Decimals:      %d\n", r.token.Decimals)
	fmt.Fprintf(&buf, "Owners:        %d\n", r.token.Stats.OwnerCounts)
	fmt.Fprintf(&buf, "Total Balance: %g\n", r.token.Stats.TotalBalance)
	return buf.String()
}

func (r *ftGetResult) Oneliner() string {
	return fmt.Sprintf("%s (%s)", r.token.Name, r.token.Symbol)
}

func (r *ftGetResult) JSON() any { return r.token }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetFT().Token(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("token %q not found", args[0])
	}
	return &ftGetResult{token: resp.Data[0]}, nil
}
