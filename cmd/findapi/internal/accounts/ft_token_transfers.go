package accounts

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type ftTokenTransfersFlags struct {
	Height uint64 `flag:"height" info:"Block height filter"`
	Limit  int    `flag:"limit"  info:"Number of results (max 100)"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var ftTokenTransfersFlagsVal = &ftTokenTransfersFlags{}

var ftTokenTransfersCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "ft-token-transfers <address> <token>",
		Short:   "List transfers for a specific FT token for an account",
		Example: "find accounts ft-token-transfers 0x1234 A.1654653399040a61.FlowToken.Vault",
		Args:    cobra.ExactArgs(2),
	},
	Flags: ftTokenTransfersFlagsVal,
	Run:   runFTTokenTransfers,
}

type ftTokenTransfersResult struct{ transfers []flow.FTTransfer }

func (r *ftTokenTransfersResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DIRECTION\tAMOUNT\tSENDER\tRECEIVER\tHEIGHT")
	for _, t := range r.transfers {
		fmt.Fprintf(w, "%s\t%g\t%s\t%s\t%d\n",
			t.Direction, t.Amount, t.Sender, t.Receiver, t.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *ftTokenTransfersResult) Oneliner() string {
	return fmt.Sprintf("%d transfers", len(r.transfers))
}

func (r *ftTokenTransfersResult) JSON() any { return r.transfers }

func runFTTokenTransfers(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountFTTokenTransfers().Address(args[0]).Token(args[1])
	if ftTokenTransfersFlagsVal.Height > 0 {
		b = b.Height(ftTokenTransfersFlagsVal.Height)
	}
	if ftTokenTransfersFlagsVal.Limit > 0 {
		b = b.Limit(ftTokenTransfersFlagsVal.Limit)
	}
	if ftTokenTransfersFlagsVal.Offset > 0 {
		b = b.Offset(ftTokenTransfersFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftTokenTransfersResult{transfers: resp.Data}, nil
}
