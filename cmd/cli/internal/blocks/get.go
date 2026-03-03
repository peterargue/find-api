package blocks

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "get <height>",
		Short:   "Get a block by height",
		Example: "find blocks get 12345678",
		Args:    cobra.ExactArgs(1),
	},
	Run: runGet,
}

type blockResult struct {
	block flow.Block
}

func (r *blockResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Height:       %d\n", r.block.Height)
	fmt.Fprintf(&buf, "ID:           %s\n", r.block.ID)
	fmt.Fprintf(&buf, "Timestamp:    %s\n", r.block.Timestamp)
	fmt.Fprintf(&buf, "Transactions: %d\n", r.block.Tx)
	fmt.Fprintf(&buf, "Fees:         %g\n", r.block.Fees)
	fmt.Fprintf(&buf, "Gas Used:     %d\n", r.block.TotalGasUsed)
	fmt.Fprintf(&buf, "EVM Txs:      %d\n", r.block.EvmTxCount)
	fmt.Fprintf(&buf, "Surge Factor: %g\n", r.block.SurgeFactor)
	return buf.String()
}

func (r *blockResult) Oneliner() string {
	return fmt.Sprintf("%d %s", r.block.Height, r.block.ID)
}

func (r *blockResult) JSON() any {
	return map[string]any{
		"height":         r.block.Height,
		"id":             r.block.ID,
		"timestamp":      r.block.Timestamp,
		"tx":             r.block.Tx,
		"fees":           r.block.Fees,
		"total_gas_used": r.block.TotalGasUsed,
		"evm_tx_count":   r.block.EvmTxCount,
		"surge_factor":   r.block.SurgeFactor,
	}
}

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid height %q: %w", args[0], err)
	}
	client := command.MustLoadClient()
	resp, err := client.Flow.GetBlock().Height(height).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("block %d not found", height)
	}
	return &blockResult{block: resp.Data[0]}, nil
}
