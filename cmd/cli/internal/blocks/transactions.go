package blocks

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type blockTxFlags struct {
	IncludeEvents bool `flag:"include-events" info:"Include transaction events in output"`
}

var blockTxFlagsVal = &blockTxFlags{}

var transactionsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "transactions <height>",
		Short:   "List transactions in a block",
		Example: "find blocks transactions 12345678\nfind blocks transactions 12345678 --include-events",
		Args:    cobra.ExactArgs(1),
	},
	Flags: blockTxFlagsVal,
	Run:   runBlockTransactions,
}

type blockTxResult struct {
	txs []flow.BlockTransaction
}

func (r *blockTxResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tFEE\tGAS\tEVENTS")
	for _, tx := range r.txs {
		fmt.Fprintf(w, "%s\t%s\t%g\t%d\t%d\n",
			tx.TransactionID, tx.Status, tx.Fee, tx.GasUsed, tx.EventCount)
	}
	w.Flush()
	return buf.String()
}

func (r *blockTxResult) Oneliner() string {
	return fmt.Sprintf("%d transactions", len(r.txs))
}

func (r *blockTxResult) JSON() any { return r.txs }

func runBlockTransactions(args []string, flags *command.GlobalFlags) (command.Result, error) {
	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid height %q: %w", args[0], err)
	}
	client := command.MustLoadClient()
	b := client.Flow.GetBlockTransactions().Height(height)
	if blockTxFlagsVal.IncludeEvents {
		b = b.IncludeEvents(true)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &blockTxResult{txs: resp.Data}, nil
}
