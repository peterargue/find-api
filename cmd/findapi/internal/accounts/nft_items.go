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

type nftItemsFlags struct {
	Limit     int    `flag:"limit"      info:"Number of results (max 100)"`
	Offset    int    `flag:"offset"     info:"Pagination offset"`
	ValidOnly bool   `flag:"valid-only" info:"Return only valid NFTs"`
	SortBy    string `flag:"sort-by"    info:"Sort order (asc or desc)"`
}

var nftItemsFlagsVal = &nftItemsFlags{}

var nftItemsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "nft-items <address> <nft-type>",
		Short:   "List NFTs of a specific type for an account",
		Example: "find accounts nft-items 0x1234 A.1234.SomeNFT",
		Args:    cobra.ExactArgs(2),
	},
	Flags: nftItemsFlagsVal,
	Run:   runNFTItems,
}

type nftItemsResult struct{ nfts []flow.AccountNFT }

func (r *nftItemsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNFT_ID\tNAME\tVALID")
	for _, n := range r.nfts {
		fmt.Fprintf(w, "%s	%d	%s	%v\n", n.ID, n.NFTId, n.Name, n.Valid)
	}
	w.Flush()
	return buf.String()
}

func (r *nftItemsResult) Oneliner() string { return fmt.Sprintf("%d NFTs", len(r.nfts)) }
func (r *nftItemsResult) JSON() any        { return r.nfts }

func runNFTItems(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountNFTs().Address(args[0]).NFTType(args[1])
	if nftItemsFlagsVal.Limit > 0 {
		b = b.Limit(nftItemsFlagsVal.Limit)
	}
	if nftItemsFlagsVal.Offset > 0 {
		b = b.Offset(nftItemsFlagsVal.Offset)
	}
	if nftItemsFlagsVal.ValidOnly {
		b = b.ValidOnly(true)
	}
	if nftItemsFlagsVal.SortBy != "" {
		b = b.SortBy(nftItemsFlagsVal.SortBy)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &nftItemsResult{nfts: resp.Data}, nil
}
