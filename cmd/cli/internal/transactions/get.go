package transactions

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "get <id>",
		Short: "Get transaction details including script and events",
		Args:  cobra.ExactArgs(1),
	},
	Run: runGet,
}

type txGetResult struct {
	tx flow.TransactionDetails
}

func (r *txGetResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ID:               %s\n", r.tx.ID)
	fmt.Fprintf(&buf, "Status:           %s\n", r.tx.Status)
	fmt.Fprintf(&buf, "Fee:              %g\n", r.tx.Fee)
	fmt.Fprintf(&buf, "Gas Used:         %d\n", r.tx.GasUsed)
	fmt.Fprintf(&buf, "Payer:            %s\n", r.tx.Payer)
	fmt.Fprintf(&buf, "Proposer:         %s\n", r.tx.Proposer)
	fmt.Fprintf(&buf, "Block Height:     %d\n", r.tx.BlockHeight)
	fmt.Fprintf(&buf, "Block ID:         %s\n", r.tx.BlockID)
	fmt.Fprintf(&buf, "Timestamp:        %s\n", r.tx.Timestamp)
	fmt.Fprintf(&buf, "Execution Effort: %g\n", r.tx.ExecutionEffort)
	fmt.Fprintf(&buf, "Surge Factor:     %g\n", r.tx.SurgeFactor)
	fmt.Fprintf(&buf, "Script:\n%s\n", r.tx.Script)
	return buf.String()
}

func (r *txGetResult) Oneliner() string {
	return fmt.Sprintf("%s %s", r.tx.ID, r.tx.Status)
}

func (r *txGetResult) JSON() any { return r.tx }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetTransaction().
		ID(args[0]).
		IncludeEvents(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("transaction not found")
	}
	return &txGetResult{tx: resp.Data[0]}, nil
}
