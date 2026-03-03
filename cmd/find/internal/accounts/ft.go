package accounts

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type ftFlags struct {
	Limit  int `flag:"limit"  info:"Number of results (max 100)"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var ftFlagsVal = &ftFlags{}

var ftCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "ft <address>",
		Short:   "List FT collections for an account",
		Example: "find accounts ft 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Flags: ftFlagsVal,
	Run:   runFT,
}

type ftResult struct{ collections []flow.AccountFTCollection }

func (r *ftResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOKEN\tBALANCE\tPATH")
	for _, c := range r.collections {
		fmt.Fprintf(w, "%s\t%s\t%s\n", c.Token, c.Balance, c.Path)
	}
	w.Flush()
	return buf.String()
}

func (r *ftResult) Oneliner() string {
	return fmt.Sprintf("%d FT collections", len(r.collections))
}

func (r *ftResult) JSON() any { return r.collections }

func runFT(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountFTs().Address(args[0])
	if ftFlagsVal.Limit > 0 {
		b = b.Limit(ftFlagsVal.Limit)
	}
	if ftFlagsVal.Offset > 0 {
		b = b.Offset(ftFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftResult{collections: resp.Data}, nil
}
