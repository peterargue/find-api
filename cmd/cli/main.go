package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:              "find",
		Short:            "CLI for the FindLabs API (api.find.xyz)",
		TraverseChildren: true,
	}

	cmd.SilenceUsage = true
	cmd.SilenceErrors = true

	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
