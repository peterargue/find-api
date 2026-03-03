package contracts

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "contracts", Short: "Query Flow contracts", TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	byIdentifierCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
}
