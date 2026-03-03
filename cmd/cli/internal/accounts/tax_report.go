package accounts

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type taxReportFlags struct {
	Height uint64 `flag:"height" info:"Block height filter"`
	Limit  int    `flag:"limit"  info:"Number of results (max 100)"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var taxReportFlagsVal = &taxReportFlags{}

var taxReportCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "tax-report <address>",
		Short:   "Get tax report for an account",
		Example: "find accounts tax-report 0x1234567890abcdef",
		Args:    cobra.ExactArgs(1),
	},
	Flags: taxReportFlagsVal,
	Run:   runTaxReport,
}

type taxReportResult struct{ entries []flow.TaxReportEntry }

func (r *taxReportResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TIME\tDIRECTION\tTYPE\tTOKEN\tAMOUNT\tFEE")
	for _, e := range r.entries {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%g\t%g\n",
			e.Time, e.Direction, e.Type, e.Token, e.Amount, e.Fee)
	}
	w.Flush()
	return buf.String()
}

func (r *taxReportResult) Oneliner() string { return fmt.Sprintf("%d entries", len(r.entries)) }
func (r *taxReportResult) JSON() any        { return r.entries }

func runTaxReport(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetAccountTaxReport().Address(args[0])
	if taxReportFlagsVal.Height > 0 {
		b = b.Height(taxReportFlagsVal.Height)
	}
	if taxReportFlagsVal.Limit > 0 {
		b = b.Limit(taxReportFlagsVal.Limit)
	}
	if taxReportFlagsVal.Offset > 0 {
		b = b.Offset(taxReportFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &taxReportResult{entries: resp.Data}, nil
}
