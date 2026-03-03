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

type accountTxFlags struct {
	Height        uint64 `flag:"height"         info:"Block height filter"`
	Limit         int    `flag:"limit"          info:"Number of results (max 100)"`
	Offset        int    `flag:"offset"         info:"Pagination offset"`
	IncludeEvents bool   `flag:"include-events" info:"Include events in response"`
	From          string `flag:"from"           info:"Start timestamp filter (ISO 8601)"`
	To            string `flag:"to"             info:"End timestamp filter (ISO 8601)"`
}

var accountTxFlagsVal = &accountTxFlags{}

var txCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "transactions <address>",
		Short:   "List transactions for an account",
		Example: "find accounts transactions 0x1234567890abcdef --limit 20",
		Args:    cobra.ExactArgs(1),
	},
	Flags: accountTxFlagsVal,
	Run:   runAccountTransactions,
}

type accountTxResult struct{ txs []flow.AccountTransaction }

func (r *accountTxResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tFEE\tGAS\tHEIGHT\tTIMESTAMP")
	for _, tx := range r.txs {
		fmt.Fprintf(w, "%s\t%s\t%g\t%d\t%d\t%s\n",
			tx.TransactionID, tx.Status, tx.Fee, tx.GasUsed, tx.BlockHeight, tx.Timestamp)
	}
	w.Flush()
	return buf.String()
}

func (r *accountTxResult) Oneliner() string { return fmt.Sprintf("%d transactions", len(r.txs)) }
func (r *accountTxResult) JSON() any        { return r.txs }

func runAccountTransactions(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountTransactions().Address(args[0])
	if accountTxFlagsVal.Height > 0 {
		b = b.Height(accountTxFlagsVal.Height)
	}
	if accountTxFlagsVal.Limit > 0 {
		b = b.Limit(accountTxFlagsVal.Limit)
	}
	if accountTxFlagsVal.Offset > 0 {
		b = b.Offset(accountTxFlagsVal.Offset)
	}
	if accountTxFlagsVal.IncludeEvents {
		b = b.IncludeEvents(true)
	}
	if accountTxFlagsVal.From != "" {
		b = b.From(accountTxFlagsVal.From)
	}
	if accountTxFlagsVal.To != "" {
		b = b.To(accountTxFlagsVal.To)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &accountTxResult{txs: resp.Data}, nil
}
