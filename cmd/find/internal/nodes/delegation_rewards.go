package nodes

import (
	"bytes"
	"context"
	"fmt"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type delegationRewardsFlags struct {
	Address string `flag:"address" info:"Filter by delegator address"`
	Limit   int    `flag:"limit"   info:"Number of rewards to return"`
	Offset  int    `flag:"offset"  info:"Pagination offset"`
	SortBy  string `flag:"sort-by" info:"Sort field (timestamp, amount)"`
}

var delegationRewardsFlagsVal = &delegationRewardsFlags{}

var delegationRewardsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "delegation-rewards <node-id>",
		Short: "List delegation rewards for a node",
		Args:  cobra.ExactArgs(1),
	},
	Flags: delegationRewardsFlagsVal,
	Run:   runDelegationRewards,
}

type delegationRewardsResult struct {
	rewards []flow.DelegationReward
}

func (r *delegationRewardsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ADDRESS\tDELEGATOR_ID\tAMOUNT\tHEIGHT\tTIMESTAMP")
	for _, rw := range r.rewards {
		fmt.Fprintf(w, "%s\t%s\t%g\t%d\t%s\n", rw.Address, rw.DelegatorID, rw.Amount, rw.BlockHeight, rw.Timestamp)
	}
	w.Flush()
	return buf.String()
}

func (r *delegationRewardsResult) Oneliner() string {
	return fmt.Sprintf("%d rewards", len(r.rewards))
}

func (r *delegationRewardsResult) JSON() any { return r.rewards }

func runDelegationRewards(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	b := client.Flow.GetNodeDelegationRewards().NodeID(args[0])
	if delegationRewardsFlagsVal.Address != "" {
		b = b.Address(delegationRewardsFlagsVal.Address)
	}
	if delegationRewardsFlagsVal.Limit > 0 {
		b = b.Limit(delegationRewardsFlagsVal.Limit)
	}
	if delegationRewardsFlagsVal.Offset > 0 {
		b = b.Offset(delegationRewardsFlagsVal.Offset)
	}
	if delegationRewardsFlagsVal.SortBy != "" {
		b = b.SortBy(delegationRewardsFlagsVal.SortBy)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &delegationRewardsResult{rewards: resp.Data}, nil
}
