package nft

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "get <type>",
		Short: "Get NFT collection details",
		Args:  cobra.ExactArgs(1),
	},
	Run: runGet,
}

type nftGetResult struct {
	collection flow.NFTCollectionDetails
}

func (r *nftGetResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Name:        %s\n", r.collection.Name)
	fmt.Fprintf(&buf, "Type:        %s\n", r.collection.NFTType)
	fmt.Fprintf(&buf, "Contract:    %s\n", r.collection.ContractName)
	fmt.Fprintf(&buf, "Address:     %s\n", r.collection.Address)
	fmt.Fprintf(&buf, "Holders:     %d\n", r.collection.HolderCount)
	fmt.Fprintf(&buf, "Items:       %d\n", r.collection.ItemCount)
	fmt.Fprintf(&buf, "Description: %s\n", r.collection.Description)
	return buf.String()
}

func (r *nftGetResult) Oneliner() string {
	return fmt.Sprintf("%s holders=%d items=%d", r.collection.Name, r.collection.HolderCount, r.collection.ItemCount)
}

func (r *nftGetResult) JSON() any { return r.collection }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetNFTCollection().NFTType(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("NFT collection %q not found", args[0])
	}
	return &nftGetResult{collection: resp.Data[0]}, nil
}
