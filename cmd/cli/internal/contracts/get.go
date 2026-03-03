package contracts

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
		Use:   "get <identifier> <id>",
		Short: "Get a specific contract by identifier and ID",
		Args:  cobra.ExactArgs(2),
	},
	Run: runGet,
}

type contractResult struct {
	contract flow.Contract
}

func (r *contractResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "identifier:    %s\n", r.contract.Identifier)
	fmt.Fprintf(&buf, "contract_name: %s\n", r.contract.ContractName)
	fmt.Fprintf(&buf, "address:       %s\n", r.contract.Address)
	fmt.Fprintf(&buf, "height:        %d\n", r.contract.BlockHeight)
	return buf.String()
}

func (r *contractResult) Oneliner() string {
	return fmt.Sprintf("%s %s", r.contract.Identifier, r.contract.Address)
}

func (r *contractResult) JSON() any { return r.contract }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetContract().Identifier(args[0]).ID(args[1]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("contract not found")
	}
	return &contractResult{contract: resp.Data[0]}, nil
}
