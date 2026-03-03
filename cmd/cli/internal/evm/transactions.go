package evm

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type transactionsFlags struct {
	Height uint64 `flag:"height" info:"Block height filter"`
	Limit  int    `flag:"limit"  info:"Number of transactions to return"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var transactionsFlagsVal = &transactionsFlags{}

var transactionsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "transactions",
		Short: "List EVM transactions",
	},
	Flags: transactionsFlagsVal,
	Run:   runTransactions,
}

type evmTransactionsResult struct {
	txs []flow.EvmTransaction
}

func (r *evmTransactionsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HASH\tFROM\tTO\tSTATUS\tBLOCK\tGAS\tVALUE")
	for _, tx := range r.txs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n", tx.Hash, tx.From, tx.To, tx.Status, tx.BlockNumber, tx.GasUsed, tx.Value)
	}
	w.Flush()
	return buf.String()
}

func (r *evmTransactionsResult) Oneliner() string {
	return fmt.Sprintf("%d transactions", len(r.txs))
}

func (r *evmTransactionsResult) JSON() any { return r.txs }

func runTransactions(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetEvmTransactions()
	if transactionsFlagsVal.Height > 0 {
		b = b.Height(transactionsFlagsVal.Height)
	}
	if transactionsFlagsVal.Limit > 0 {
		b = b.Limit(transactionsFlagsVal.Limit)
	}
	if transactionsFlagsVal.Offset > 0 {
		b = b.Offset(transactionsFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &evmTransactionsResult{txs: resp.Data}, nil
}
