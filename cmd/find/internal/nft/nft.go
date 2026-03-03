package nft

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "nft", Short: "Query NFT collections", TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	transfersCmd.AddToParent(Cmd)
	holdingsCmd.AddToParent(Cmd)
	itemCmd.AddToParent(Cmd)
}
