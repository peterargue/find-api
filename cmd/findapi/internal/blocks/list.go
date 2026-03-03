package blocks

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type listFlags struct {
	Height uint64 `flag:"height" info:"Starting block height (returns blocks at or below this height, descending)"`
	Limit  int    `flag:"limit"  info:"Number of blocks to return (max 100, default 25)"`
	Offset int    `flag:"offset" info:"Pagination offset"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "list",
		Short:   "List recent blocks",
		Example: "find blocks list\nfind blocks list --height 12345678 --limit 10",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type blocksResult struct {
	blocks []flow.Block
}

func (r *blocksResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HEIGHT\tID\tTIMESTAMP\tTXS\tFEES")
	for _, b := range r.blocks {
		fmt.Fprintf(w, "%d\t%s\t%s\t%d\t%g\n", b.Height, b.ID, b.Timestamp, b.Tx, b.Fees)
	}
	w.Flush()
	return buf.String()
}

func (r *blocksResult) Oneliner() string {
	if len(r.blocks) == 0 {
		return "no blocks found"
	}
	b := r.blocks[0]
	return fmt.Sprintf("%d %s", b.Height, b.ID)
}

func (r *blocksResult) JSON() any { return r.blocks }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetBlocks()
	if listFlagsVal.Height > 0 {
		b = b.Height(listFlagsVal.Height)
	}
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &blocksResult{blocks: resp.Data}, nil
}
