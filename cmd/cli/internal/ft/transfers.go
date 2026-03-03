package ft

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type transfersFlags struct {
	Token  string `flag:"token"   info:"Token identifier filter"`
	TxHash string `flag:"tx-hash" info:"Transaction hash filter"`
	Height uint64 `flag:"height"  info:"Block height filter"`
	Limit  int    `flag:"limit"   info:"Number of transfers to return"`
	Offset int    `flag:"offset"  info:"Pagination offset"`
}

var transfersFlagsVal = &transfersFlags{}

var transfersCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "transfers",
		Short: "List fungible token transfers",
	},
	Flags: transfersFlagsVal,
	Run:   runTransfers,
}

type ftTransfersResult struct {
	transfers []flow.FTTransfer
}

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

func runTransfers(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetFTTransfers()
	if transfersFlagsVal.Token != "" {
		b = b.Token(transfersFlagsVal.Token)
	}
	if transfersFlagsVal.TxHash != "" {
		b = b.TransactionHash(transfersFlagsVal.TxHash)
	}
	if transfersFlagsVal.Height > 0 {
		b = b.Height(transfersFlagsVal.Height)
	}
	if transfersFlagsVal.Limit > 0 {
		b = b.Limit(transfersFlagsVal.Limit)
	}
	if transfersFlagsVal.Offset > 0 {
		b = b.Offset(transfersFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftTransfersResult{transfers: resp.Data}, nil
}
