package ft

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
	Height uint64 `flag:"height" info:"Block height filter"`
	Limit  int    `flag:"limit"  info:"Number of tokens to return"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "list",
		Short: "List fungible tokens",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type ftListResult struct {
	tokens []flow.FungibleToken
}

func (r *ftListResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOKEN\tNAME\tSYMBOL\tSUPPLY\tDECIMALS")
	for _, t := range r.tokens {
		fmt.Fprintf(w, "%s\t%s\t%s\t%g\t%d\n", t.Token, t.Name, t.Symbol, t.CirculatingSupply, t.Decimals)
	}
	w.Flush()
	return buf.String()
}

func (r *ftListResult) Oneliner() string {
	return fmt.Sprintf("%d tokens", len(r.tokens))
}

func (r *ftListResult) JSON() any { return r.tokens }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetFTs()
	if listFlagsVal.Height > 0 {
		b = b.Height(listFlagsVal.Height)
	}
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
	return &ftListResult{tokens: resp.Data}, nil
}
