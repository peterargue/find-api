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

type listFlags struct {
	Limit  int `flag:"limit"  info:"Number of contracts to return"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "list",
		Short: "List contracts",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type contractsResult struct {
	contracts []flow.Contract
}

func (r *contractsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "IDENTIFIER\tCONTRACT_NAME\tADDRESS\tHEIGHT")
	for _, c := range r.contracts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%d\n", c.Identifier, c.ContractName, c.Address, c.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *contractsResult) Oneliner() string {
	return fmt.Sprintf("%d contracts", len(r.contracts))
}

func (r *contractsResult) JSON() any { return r.contracts }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetContracts()
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &contractsResult{contracts: resp.Data}, nil
}
