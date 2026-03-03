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

type nftFlags struct {
	Limit  int `flag:"limit"  info:"Number of results (max 100)"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var nftFlagsVal = &nftFlags{}

var nftCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "nft <address>",
		Short:   "List NFT collections for an account",
		Example: "find accounts nft 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Flags: nftFlagsVal,
	Run:   runNFT,
}

type nftResult struct{ collections []flow.AccountNFTCollection }

func (r *nftResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tNFT_TYPE\tCOUNT")
	for _, c := range r.collections {
		fmt.Fprintf(w, "%s\t%s\t%d\n", c.Name, c.NFTType, c.NFTCount)
	}
	w.Flush()
	return buf.String()
}

func (r *nftResult) Oneliner() string { return fmt.Sprintf("%d NFT collections", len(r.collections)) }
func (r *nftResult) JSON() any        { return r.collections }

func runNFT(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountNFTCollections().Address(args[0])
	if nftFlagsVal.Limit > 0 {
		b = b.Limit(nftFlagsVal.Limit)
	}
	if nftFlagsVal.Offset > 0 {
		b = b.Offset(nftFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &nftResult{collections: resp.Data}, nil
}
