package command

import "github.com/spf13/cobra"

// GlobalFlags are available on every command via persistent flags on the root.
type GlobalFlags struct {
	Format string
	Filter string
	Save   string
	Log    string
}

// Flags holds the parsed global flags for the current invocation.
var Flags GlobalFlags

// InitFlags registers global flags on the root command.
func InitFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&Flags.Format, "format", "text", "Output format: text, json, inline")
	cmd.PersistentFlags().StringVar(&Flags.Filter, "filter", "", "Filter output by field name")
	cmd.PersistentFlags().StringVar(&Flags.Save, "save", "", "Save output to file")
	cmd.PersistentFlags().StringVar(&Flags.Log, "log", "info", "Log level: debug, info, error, none")
}
