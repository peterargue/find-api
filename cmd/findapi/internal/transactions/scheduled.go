package transactions

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type scheduledFlags struct {
	Completed bool   `flag:"completed" info:"Filter by completed status"`
	Owner     string `flag:"owner"     info:"Owner address filter"`
	Status    string `flag:"status"    info:"Status filter"`
	Limit     int    `flag:"limit"     info:"Number of scheduled transactions to return"`
	Offset    int    `flag:"offset"    info:"Pagination offset"`
}

var scheduledFlagsVal = &scheduledFlags{}

var scheduledCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "scheduled",
		Short: "List scheduled transactions",
	},
	Flags: scheduledFlagsVal,
	Run:   runScheduled,
}

type scheduledResult struct {
	txs []flow.ScheduledTransaction
}

func (r *scheduledResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tSTATUS\tOWNER\tPRIORITY\tCOMPLETED\tSCHEDULED_AT")
	for _, t := range r.txs {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%t\t%s\n",
			t.ID, t.Status, t.Owner, t.Priority, t.IsCompleted, t.ScheduledAt)
	}
	w.Flush()
	return buf.String()
}

func (r *scheduledResult) Oneliner() string {
	return fmt.Sprintf("%d scheduled transactions", len(r.txs))
}

func (r *scheduledResult) JSON() any { return r.txs }

func runScheduled(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetScheduledTransactions()
	if scheduledFlagsVal.Completed {
		b = b.Completed(true)
	}
	if scheduledFlagsVal.Owner != "" {
		b = b.Owner(scheduledFlagsVal.Owner)
	}
	if scheduledFlagsVal.Status != "" {
		b = b.Status(scheduledFlagsVal.Status)
	}
	if scheduledFlagsVal.Limit > 0 {
		b = b.Limit(scheduledFlagsVal.Limit)
	}
	if scheduledFlagsVal.Offset > 0 {
		b = b.Offset(scheduledFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &scheduledResult{txs: resp.Data}, nil
}
