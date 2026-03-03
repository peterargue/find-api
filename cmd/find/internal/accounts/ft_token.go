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

type ftTokenFlags struct {
	Limit  int `flag:"limit"  info:"Number of results (max 100)"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var ftTokenFlagsVal = &ftTokenFlags{}

var ftTokenCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "ft-token <address> <token>",
		Short:   "Get vault info for a specific FT token for an account",
		Example: "find accounts ft-token 0x1234 A.1654653399040a61.FlowToken.Vault",
		Args:    cobra.ExactArgs(2),
	},
	Flags: ftTokenFlagsVal,
	Run:   runFTToken,
}

type ftTokenResult struct{ vaults []flow.Vault }

func (r *ftTokenResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TOKEN\tBALANCE\tPATH\tHEIGHT")
	for _, v := range r.vaults {
		fmt.Fprintf(w, "%s\t%g\t%s\t%d\n", v.Token, v.Balance, v.Path, v.BlockHeight)
	}
	w.Flush()
	return buf.String()
}

func (r *ftTokenResult) Oneliner() string {
	if len(r.vaults) == 0 {
		return "no vault found"
	}
	return fmt.Sprintf("%s balance=%g", r.vaults[0].Token, r.vaults[0].Balance)
}

func (r *ftTokenResult) JSON() any { return r.vaults }

func runFTToken(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountFTToken().Address(args[0]).Token(args[1])
	if ftTokenFlagsVal.Limit > 0 {
		b = b.Limit(ftTokenFlagsVal.Limit)
	}
	if ftTokenFlagsVal.Offset > 0 {
		b = b.Offset(ftTokenFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &ftTokenResult{vaults: resp.Data}, nil
}
