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

type ftTransfersFlags struct {
	Height uint64 `flag:"height" info:"Block height filter"`
	Limit  int    `flag:"limit"  info:"Number of results (max 100)"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var ftTransfersFlagsVal = &ftTransfersFlags{}

var ftTransfersCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "ft-transfers <address>",
		Short:   "List FT transfers for an account",
		Example: "find accounts ft-transfers 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Flags: ftTransfersFlagsVal,
	Run:   runFTTransfers,
}

type ftTransfersResult struct{ transfers []flow.FTTransfer }

func (r *ftTransfersResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DIRECTION\tAMOUNT\tTOKEN\tSENDER\tRECEIVER\tHEIGHT")
	for _, t := range r.transfers {
		fmt.Fprintf(w, "%s\t%g\t%s\t%s\t%s\t%d\n",
			t.Direction, t.Amount, t.Token.Token, t.Sender, t.Receiver, t.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *ftTransfersResult) Oneliner() string {
	return fmt.Sprintf("%d transfers", len(r.transfers))
}

func (r *ftTransfersResult) JSON() any { return r.transfers }

func runFTTransfers(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountFTTransfers().Address(args[0])
	if ftTransfersFlagsVal.Height > 0 {
		b = b.Height(ftTransfersFlagsVal.Height)
	}
	if ftTransfersFlagsVal.Limit > 0 {
		b = b.Limit(ftTransfersFlagsVal.Limit)
	}
	if ftTransfersFlagsVal.Offset > 0 {
		b = b.Offset(ftTransfersFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftTransfersResult{transfers: resp.Data}, nil
}
