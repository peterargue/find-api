package accounts

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type ftHoldingsFlags struct {
	Limit  int `flag:"limit"  info:"Number of results (max 100)"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var ftHoldingsFlagsVal = &ftHoldingsFlags{}

var ftHoldingsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "ft-holdings <address>",
		Short:   "List FT holdings with statistics for an account",
		Example: "find accounts ft-holdings 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Flags: ftHoldingsFlagsVal,
	Run:   runFTHoldings,
}

type ftHoldingsResult struct{ holdings []flow.FTHolding }

func (r *ftHoldingsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOKEN\tBALANCE\tPERCENTAGE")
	for _, h := range r.holdings {
		fmt.Fprintf(w, "%s\t%g\t%.4f%%\n", h.Token, h.Balance, h.Percentage)
	}
	w.Flush()
	return buf.String()
}

func (r *ftHoldingsResult) Oneliner() string { return fmt.Sprintf("%d holdings", len(r.holdings)) }
func (r *ftHoldingsResult) JSON() any        { return r.holdings }

func runFTHoldings(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountFTHoldings().Address(args[0])
	if ftHoldingsFlagsVal.Limit > 0 {
		b = b.Limit(ftHoldingsFlagsVal.Limit)
	}
	if ftHoldingsFlagsVal.Offset > 0 {
		b = b.Offset(ftHoldingsFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftHoldingsResult{holdings: resp.Data}, nil
}
