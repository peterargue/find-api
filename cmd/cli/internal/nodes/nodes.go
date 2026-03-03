package nodes

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "nodes", Short: "Query Flow staking nodes", TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	delegationRewardsCmd.AddToParent(Cmd)
}
