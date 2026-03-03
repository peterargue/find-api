package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/peterargue/find-api/cmd/find/internal/accounts"
	"github.com/peterargue/find-api/cmd/find/internal/auth"
	"github.com/peterargue/find-api/cmd/find/internal/blocks"
	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/cmd/find/internal/contracts"
	"github.com/peterargue/find-api/cmd/find/internal/evm"
	"github.com/peterargue/find-api/cmd/find/internal/ft"
	"github.com/peterargue/find-api/cmd/find/internal/nft"
	"github.com/peterargue/find-api/cmd/find/internal/nodes"
	"github.com/peterargue/find-api/cmd/find/internal/transactions"
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
