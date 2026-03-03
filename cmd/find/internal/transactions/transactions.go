package transactions

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "transactions", Short: "Query Flow transactions", TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	scheduledCmd.AddToParent(Cmd)
}
