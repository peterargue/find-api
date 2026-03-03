package contracts

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type byIdentifierFlags struct {
	Limit  int `flag:"limit"  info:"Number of contracts to return"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var byIdentifierFlagsVal = &byIdentifierFlags{}

var byIdentifierCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "by-identifier <identifier>",
		Short: "List contracts by identifier",
		Args:  cobra.ExactArgs(1),
	},
	Flags: byIdentifierFlagsVal,
	Run:   runByIdentifier,
}

type contractsByIdentifierResult struct {
	contracts []flow.Contract
}

func (r *contractsByIdentifierResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "IDENTIFIER\tCONTRACT_NAME\tADDRESS\tHEIGHT")
	for _, c := range r.contracts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", c.Identifier, c.ContractName, c.Address, c.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *contractsByIdentifierResult) Oneliner() string {
	return fmt.Sprintf("%d contracts", len(r.contracts))
}

func (r *contractsByIdentifierResult) JSON() any { return r.contracts }

func runByIdentifier(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetContractsByIdentifier().Identifier(args[0])
	if byIdentifierFlagsVal.Limit > 0 {
		b = b.Limit(byIdentifierFlagsVal.Limit)
	}
	if byIdentifierFlagsVal.Offset > 0 {
		b = b.Offset(byIdentifierFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &contractsByIdentifierResult{contracts: resp.Data}, nil
}
