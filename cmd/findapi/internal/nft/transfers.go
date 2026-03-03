package nft

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type transfersFlags struct {
	Address string `flag:"address"  info:"Address filter"`
	NFTType string `flag:"nft-type" info:"NFT type filter"`
	Height  uint64 `flag:"height"   info:"Block height filter"`
	Limit   int    `flag:"limit"    info:"Number of transfers to return"`
	Offset  int    `flag:"offset"   info:"Pagination offset"`
}

var transfersFlagsVal = &transfersFlags{}

var transfersCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "transfers",
		Short: "List NFT transfers",
	},
	Flags: transfersFlagsVal,
	Run:   runTransfers,
}

type nftTransfersResult struct {
	transfers []flow.NFTTransfer
}

func (r *nftTransfersResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "DIRECTION\tNFT_ID\tNFT_TYPE\tSENDER\tRECEIVER\tHEIGHT")
	for _, t := range r.transfers {
		fmt.Fprintf(w, "%s\t%d\t%s\t%s\t%s\t%d\n",
			t.Direction, t.NFTId, t.NFTType, t.Sender, t.Receiver, t.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *nftTransfersResult) Oneliner() string {
	return fmt.Sprintf("%d transfers", len(r.transfers))
}

func (r *nftTransfersResult) JSON() any { return r.transfers }

func runTransfers(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetNFTTransfers()
	if transfersFlagsVal.Address != "" {
		b = b.Address(transfersFlagsVal.Address)
	}
	if transfersFlagsVal.NFTType != "" {
		b = b.NFTType(transfersFlagsVal.NFTType)
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
	return &nftTransfersResult{transfers: resp.Data}, nil
}
