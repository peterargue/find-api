package evm

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var transactionCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "transaction <hash>",
		Short: "Get an EVM transaction by hash",
		Args:  cobra.ExactArgs(1),
	},
	Run: runTransaction,
}

type evmTransactionResult struct {
	tx flow.EvmTransaction
}

func (r *evmTransactionResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "hash:         %s\n", r.tx.Hash)
	fmt.Fprintf(&buf, "from:         %s\n", r.tx.From)
	fmt.Fprintf(&buf, "to:           %s\n", r.tx.To)
	fmt.Fprintf(&buf, "status:       %s\n", r.tx.Status)
	fmt.Fprintf(&buf, "block_number: %d\n", r.tx.BlockNumber)
	fmt.Fprintf(&buf, "gas_used:     %s\n", r.tx.GasUsed)
	fmt.Fprintf(&buf, "gas_limit:    %s\n", r.tx.GasLimit)
	fmt.Fprintf(&buf, "value:        %s\n", r.tx.Value)
	fmt.Fprintf(&buf, "nonce:        %d\n", r.tx.Nonce)
	fmt.Fprintf(&buf, "timestamp:    %s\n", r.tx.Timestamp)
	return buf.String()
}

func (r *evmTransactionResult) Oneliner() string {
	return fmt.Sprintf("%s %s", r.tx.Hash, r.tx.Status)
}

func (r *evmTransactionResult) JSON() any { return r.tx }

func runTransaction(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	tx, err := client.Flow.GetEvmTransaction().Hash(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, fmt.Errorf("EVM transaction not found")
	}
	return &evmTransactionResult{tx: *tx}, nil
}
