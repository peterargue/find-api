package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/peterargue/find-api/cmd/cli/internal/accounts"
	"github.com/peterargue/find-api/cmd/cli/internal/auth"
	"github.com/peterargue/find-api/cmd/cli/internal/blocks"
	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/cmd/cli/internal/contracts"
	"github.com/peterargue/find-api/cmd/cli/internal/evm"
	"github.com/peterargue/find-api/cmd/cli/internal/ft"
	"github.com/peterargue/find-api/cmd/cli/internal/nft"
	"github.com/peterargue/find-api/cmd/cli/internal/nodes"
	"github.com/peterargue/find-api/cmd/cli/internal/transactions"
)

func main() {
	cmd := &cobra.Command{
		Use:              "find",
		Short:            "CLI for the FindLabs API (api.find.xyz)",
		TraverseChildren: true,
	}

	cmd.AddCommand(auth.Cmd)
	cmd.AddCommand(blocks.Cmd)
	cmd.AddCommand(accounts.Cmd)
	cmd.AddCommand(ft.Cmd)
	cmd.AddCommand(nft.Cmd)
	cmd.AddCommand(transactions.Cmd)
	cmd.AddCommand(nodes.Cmd)
	cmd.AddCommand(contracts.Cmd)
	cmd.AddCommand(evm.Cmd)

	command.InitFlags(cmd)

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
