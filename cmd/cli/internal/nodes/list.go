package nodes

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/cli/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type listFlags struct {
	Height uint64 `flag:"height"  info:"Block height filter"`
	Limit  int    `flag:"limit"   info:"Number of nodes to return"`
	Offset int    `flag:"offset"  info:"Pagination offset"`
	RoleID string `flag:"role-id" info:"Filter by role ID (1=collection, 2=consensus, 3=execution, 4=verification, 5=access)"`
	SortBy string `flag:"sort-by" info:"Sort field (tokens_staked, delegators)"`
}

var listFlagsVal = &listFlags{}

var listCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "list",
		Short: "List staking nodes",
	},
	Flags: listFlagsVal,
	Run:   runList,
}

type nodesResult struct {
	nodes []flow.Node
}

func (r *nodesResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "NODE_ID\tNAME\tROLE\tTOKENS_STAKED\tDELEGATORS")
	for _, n := range r.nodes {
		fmt.Fprintf(w, "%s\t%s\t%s\t%g\t%d\n", n.NodeID, n.Name, n.Role, n.TokensStaked, n.Delegators)
	}
	w.Flush()
	return buf.String()
}

func (r *nodesResult) Oneliner() string {
	return fmt.Sprintf("%d nodes", len(r.nodes))
}

func (r *nodesResult) JSON() any { return r.nodes }

func runList(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetNodes()
	if listFlagsVal.Height > 0 {
		b = b.Height(listFlagsVal.Height)
	}
	if listFlagsVal.Limit > 0 {
		b = b.Limit(listFlagsVal.Limit)
	}
	if listFlagsVal.Offset > 0 {
		b = b.Offset(listFlagsVal.Offset)
	}
	if listFlagsVal.RoleID != "" {
		b = b.RoleID(listFlagsVal.RoleID)
	}
	if listFlagsVal.SortBy != "" {
		b = b.SortBy(listFlagsVal.SortBy)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &nodesResult{nodes: resp.Data}, nil
}
