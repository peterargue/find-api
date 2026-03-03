package nft

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var itemCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "item <type> <id>",
		Short: "Get details for a specific NFT item",
		Args:  cobra.ExactArgs(2),
	},
	Run: runItem,
}

type nftItemResult struct {
	nft flow.NFT
}

func (r *nftItemResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "ID:        %s\n", r.nft.ID)
	fmt.Fprintf(&buf, "NFT ID:    %s\n", r.nft.NFTId)
	fmt.Fprintf(&buf, "Type:      %s\n", r.nft.NFTType)
	fmt.Fprintf(&buf, "Name:      %s\n", r.nft.Name)
	fmt.Fprintf(&buf, "Owner:     %s\n", r.nft.Owner)
	fmt.Fprintf(&buf, "Thumbnail: %s\n", r.nft.Thumbnail)
	fmt.Fprintf(&buf, "Height:    %d\n", r.nft.BlockHeight)
	return buf.String()
}

func (r *nftItemResult) Oneliner() string {
	return fmt.Sprintf("%s owner=%s", r.nft.Name, r.nft.Owner)
}

func (r *nftItemResult) JSON() any { return r.nft }

func runItem(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetNFTItem().NFTType(args[0]).ID(args[1]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("NFT not found")
	}
	return &nftItemResult{nft: resp.Data[0]}, nil
}
