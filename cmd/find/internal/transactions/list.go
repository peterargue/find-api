package transactions

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
	Height        uint64 `flag:"height"         info:"Block height filter"`
	Limit         int    `flag:"limit"          info:"Number of transactions to return"`
	Offset        int    `flag:"offset"         info:"Pagination offset"`
	Status        string `flag:"status"         info:"Status filter (e.g. SEALED, ERROR)"`
	Payer         string `flag:"payer"          info:"Payer address filter"`
	Proposer      string `flag:"proposer"       info:"Proposer address filter"`
	From          string `flag:"from"           info:"Start timestamp filter (ISO 8601)"`
	To            string `flag:"to"             info:"End timestamp filter (ISO 8601)"`
	IncludeEvents bool   `flag:"include-events" info:"Include events in response"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "list",
		Short: "List transactions",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type txListResult struct {
	txs []flow.Transaction
}

func (r *txListResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tFEE\tGAS\tHEIGHT\tTIMESTAMP")
	for _, t := range r.txs {
		fmt.Fprintf(w, "%s\t%s\t%g\t%d\t%d\t%s\n",
			t.ID, t.Status, t.Fee, t.GasUsed, t.BlockHeight, t.Timestamp)
	}
	w.Flush()
	return buf.String()
}

func (r *txListResult) Oneliner() string {
	return fmt.Sprintf("%d transactions", len(r.txs))
}

func (r *txListResult) JSON() any { return r.txs }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetTransactions()
	if listFlagsVal.Height > 0 {
		b = b.Height(listFlagsVal.Height)
	}
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	if listFlagsVal.Status != "" {
		b = b.Status(listFlagsVal.Status)
	}
	if listFlagsVal.Payer != "" {
		b = b.Payer(listFlagsVal.Payer)
	}
	if listFlagsVal.Proposer != "" {
		b = b.Proposer(listFlagsVal.Proposer)
	}
	if listFlagsVal.From != "" {
		b = b.From(listFlagsVal.From)
	}
	if listFlagsVal.To != "" {
		b = b.To(listFlagsVal.To)
	}
	if listFlagsVal.IncludeEvents {
		b = b.IncludeEvents(true)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &txListResult{txs: resp.Data}, nil
}
