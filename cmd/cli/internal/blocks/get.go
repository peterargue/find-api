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

func (r *blockResult) JSON() any { return r.block }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	ctx := context.Background()

	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%q looks like a block ID; block lookup by ID is not supported — please provide a block height instead", args[0])
	}

	resp, err := client.Flow.GetBlock().Height(height).Do(ctx)
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("block %d not found", height)
	}
	return &blockResult{block: resp.Data[0]}, nil
}
