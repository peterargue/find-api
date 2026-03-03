package ft

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type holdingsFlags struct {
	Limit  int `flag:"limit"  info:"Number of holdings to return"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var holdingsFlagsVal = &holdingsFlags{}

var holdingsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "holdings <token>",
		Short: "List fungible token holdings",
		Args:  cobra.ExactArgs(1),
	},
	Flags: holdingsFlagsVal,
	Run:   runHoldings,
}

type ftHoldingsResult struct {
	holdings []flow.FTHolding
}

func (r *ftHoldingsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ADDRESS\tBALANCE\tPERCENTAGE")
	for _, h := range r.holdings {
		fmt.Fprintf(w, "%s\t%g\t%g\n", h.Address, h.Balance, h.Percentage)
	}
	w.Flush()
	return buf.String()
}

func (r *ftHoldingsResult) Oneliner() string {
	return fmt.Sprintf("%d holders", len(r.holdings))
}

func (r *ftHoldingsResult) JSON() any { return r.holdings }

func runHoldings(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetFTHoldings().Token(args[0])
	if holdingsFlagsVal.Limit > 0 {
		b = b.Limit(holdingsFlagsVal.Limit)
	}
	if holdingsFlagsVal.Offset > 0 {
		b = b.Offset(holdingsFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftHoldingsResult{holdings: resp.Data}, nil
}
