package ft

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "ft", Short: "Query fungible tokens", TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	transfersCmd.AddToParent(Cmd)
	holdingsCmd.AddToParent(Cmd)
}
