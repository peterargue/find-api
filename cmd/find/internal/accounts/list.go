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

type listFlags struct {
	Height uint64 `flag:"height"  info:"Block height cursor"`
	Limit  int    `flag:"limit"   info:"Number of accounts (max 100)"`
	Offset int    `flag:"offset"  info:"Pagination offset"`
	SortBy string `flag:"sort-by" info:"Sort field (e.g. flow_balance)"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "list",
		Short:   "List accounts",
		Example: "find accounts list --sort-by flow_balance --limit 10",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type accountsResult struct{ accounts []flow.Account }

func (r *accountsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ADDRESS\tBALANCE\tSTORAGE\tFIND_NAME")
	for _, a := range r.accounts {
		fmt.Fprintf(w, "%s\t%g\t%g\t%s\n", a.Address, a.FlowBalance, a.FlowStorage, a.FindName)
	}
	w.Flush()
	return buf.String()
}

func (r *accountsResult) Oneliner() string {
	if len(r.accounts) == 0 {
		return "no accounts found"
	}
	return r.accounts[0].Address
}

func (r *accountsResult) JSON() any { return r.accounts }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccounts()
	if listFlagsVal.Height > 0 {
		b = b.Height(listFlagsVal.Height)
	}
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	if listFlagsVal.SortBy != "" {
		b = b.SortBy(listFlagsVal.SortBy)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &accountsResult{accounts: resp.Data}, nil
}
