package command

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"

	findapi "github.com/peterargue/find-api"
	"github.com/spf13/cobra"
)

// Result is implemented by every command's output type.
type Result interface {
	String() string
	Oneliner() string
	JSON() any
}

// RunFunc is the signature for command implementations.
type RunFunc func(args []string, flags *GlobalFlags) (Result, error)

// Command wraps a cobra.Command with typed flags and a run function.
type Command struct {
	Cmd   *cobra.Command
	Flags any
	Run   RunFunc
}

// AddToParent registers the command with a parent cobra command.
// It handles output formatting, error display, and global flags.
func (c *Command) AddToParent(parent *cobra.Command) {
	c.Cmd.RunE = func(cmd *cobra.Command, args []string) error {
		result, err := c.Run(args, &Flags)
		if err != nil {
			handleError(err)
			return nil  // handleError already printed; don't let cobra print again
		}
		if result == nil {
			return nil
		}
		formatted, err := FormatResult(result, Flags.Filter, Flags.Format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			return err
		}
		printResult(formatted, Flags.Format, Flags.Filter, Flags.Save)
		return nil
	}
	bindFlags(c)
	parent.AddCommand(c.Cmd)
}

// handleError formats errors for display, with special handling for 401.
func handleError(err error) {
	var apiErr *findapi.APIError
	if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusUnauthorized {
		fmt.Fprintln(os.Stderr, "Not authenticated. Run: find auth login")
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "Error: %s\n", err)
}

// bindFlags binds the command-specific flags struct to the cobra command
// using reflection on struct field tags.
func bindFlags(c *Command) {
	if c.Flags == nil {
		return
	}
	rv := reflect.ValueOf(c.Flags)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return
	}
	v := rv.Elem()
	if v.Kind() != reflect.Struct {
		return
	}
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		flagName := field.Tag.Get("flag")
		usage := field.Tag.Get("info")
		if flagName == "" {
			continue
		}
		fv := v.Field(i)
		switch fv.Interface().(type) {
		case string:
			c.Cmd.Flags().StringVar(fv.Addr().Interface().(*string), flagName, "", usage)
		case bool:
			c.Cmd.Flags().BoolVar(fv.Addr().Interface().(*bool), flagName, false, usage)
		case int:
			c.Cmd.Flags().IntVar(fv.Addr().Interface().(*int), flagName, 0, usage)
		case uint64:
			c.Cmd.Flags().Uint64Var(fv.Addr().Interface().(*uint64), flagName, 0, usage)
		case []string:
			c.Cmd.Flags().StringSliceVar(fv.Addr().Interface().(*[]string), flagName, nil, usage)
		}
	}
}
