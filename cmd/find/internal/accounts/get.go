package accounts

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

func formatBytes(b float64) string {
	const unit = 1000.0
	if b < unit {
		return fmt.Sprintf("%.0f B", b)
	}
	div, exp := unit, 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", b/div, "KMGTPE"[exp])
}

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "get <address>",
		Short:   "Get account details",
		Example: "find accounts get 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Run: runGet,
}

type accountResult struct{ account flow.CombinedAccountDetails }

func (r *accountResult) String() string {
	var buf bytes.Buffer
	a := r.account
	fmt.Fprintf(&buf, "Address:           %s\n", a.Address)
	fmt.Fprintf(&buf, "Flow Balance:      %s\n", strconv.FormatFloat(a.FlowBalance, 'f', -1, 64))
	fmt.Fprintf(&buf, "Storage Used:      %s\n", formatBytes(a.StorageUsed))
	fmt.Fprintf(&buf, "Storage Available: %s\n", formatBytes(a.StorageAvailable))
	if a.Find != nil && a.Find.Name != "" {
		fmt.Fprintf(&buf, "Find Name:         %s\n", a.Find.Name)
	}
	fmt.Fprintf(&buf, "Contracts:         %v\n", a.Contracts)
	fmt.Fprintf(&buf, "Keys:              %d\n", len(a.Keys))
	return buf.String()
}

func (r *accountResult) Oneliner() string {
	return fmt.Sprintf("%s balance=%g", r.account.Address, r.account.FlowBalance)
}

func (r *accountResult) JSON() any { return r.account }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetAccount().Address(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("account %s not found", args[0])
	}
	return &accountResult{account: resp.Data[0]}, nil
}
