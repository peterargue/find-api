package accounts

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:              "accounts",
	Short:            "Query Flow accounts",
	TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	ftCmd.AddToParent(Cmd)
	ftHoldingsCmd.AddToParent(Cmd)
	ftTransfersCmd.AddToParent(Cmd)
	ftTokenCmd.AddToParent(Cmd)
	ftTokenTransfersCmd.AddToParent(Cmd)
	nftCmd.AddToParent(Cmd)
	nftItemsCmd.AddToParent(Cmd)
	txCmd.AddToParent(Cmd)
	taxReportCmd.AddToParent(Cmd)
}
