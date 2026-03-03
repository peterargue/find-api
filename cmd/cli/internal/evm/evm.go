package evm

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use: "evm", Short: "Query Flow EVM layer", TraverseChildren: true,
}

func init() {
	tokensCmd.AddToParent(Cmd)
	tokenCmd.AddToParent(Cmd)
	transactionsCmd.AddToParent(Cmd)
	transactionCmd.AddToParent(Cmd)
}
