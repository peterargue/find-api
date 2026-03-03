package evm

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type tokensFlags struct {
	Type   string `flag:"type"   info:"Token type filter"`
	Name   string `flag:"name"   info:"Partial name or symbol to search for"`
	Limit  int    `flag:"limit"  info:"Number of tokens to return"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var tokensFlagsVal = &tokensFlags{}

var tokensCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "tokens",
		Short: "List EVM tokens",
	},
	Flags: tokensFlagsVal,
	Run:   runTokens,
}

type evmTokensResult struct {
	tokens []flow.EvmToken
}

func (r *evmTokensResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NAME\tSYMBOL\tTYPE\tADDRESS\tDECIMALS\tHOLDERS")
	for _, t := range r.tokens {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%d\n", t.Name, t.Symbol, t.Type, t.ContractAddressHash, t.Decimals, t.Holders)
	}
	w.Flush()
	return buf.String()
}

func (r *evmTokensResult) Oneliner() string {
	return fmt.Sprintf("%d tokens", len(r.tokens))
}

func (r *evmTokensResult) JSON() any { return r.tokens }

func runTokens(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetEvmTokens()
	if tokensFlagsVal.Type != "" {
		b = b.Type(tokensFlagsVal.Type)
	}
	if tokensFlagsVal.Name != "" {
		b = b.Name(tokensFlagsVal.Name)
	}
	if tokensFlagsVal.Limit > 0 {
		b = b.Limit(tokensFlagsVal.Limit)
	}
	if tokensFlagsVal.Offset > 0 {
		b = b.Offset(tokensFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &evmTokensResult{tokens: resp.Data}, nil
}
