package blocks

import (
	"github.com/spf13/cobra"
)

// Cmd is the "blocks" subcommand group.
var Cmd = &cobra.Command{
	Use:              "blocks",
	Short:            "Query Flow blocks",
	TraverseChildren: true,
}

func init() {
	listCmd.AddToParent(Cmd)
	getCmd.AddToParent(Cmd)
	serviceEventsCmd.AddToParent(Cmd)
	transactionsCmd.AddToParent(Cmd)
}
