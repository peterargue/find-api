package nft

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type listFlags struct {
	Limit  int `flag:"limit"  info:"Number of collections to return"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "list",
		Short: "List NFT collections",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type nftListResult struct {
	collections []flow.NFTCollection
}

func (r *nftListResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tNFT_TYPE\tCONTRACT\tADDRESS")
	for _, c := range r.collections {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", c.Name, c.NFTType, c.ContractName, c.Address)
	}
	w.Flush()
	return buf.String()
}

func (r *nftListResult) Oneliner() string {
	return fmt.Sprintf("%d collections", len(r.collections))
}

func (r *nftListResult) JSON() any { return r.collections }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetNFTCollections()
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &nftListResult{collections: resp.Data}, nil
}
