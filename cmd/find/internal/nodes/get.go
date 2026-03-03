package nodes

import (
	"bytes"
	"context"
	"fmt"

	"github.com/peterargue/find-api/cmd/find/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

var getCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:   "get <node-id>",
		Short: "Get a node by ID",
		Args:  cobra.ExactArgs(1),
	},
	Run: runGet,
}

type nodeResult struct {
	node flow.Node
}

func (r *nodeResult) String() string {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "node_id:          %s\n", r.node.NodeID)
	fmt.Fprintf(&buf, "name:             %s\n", r.node.Name)
	fmt.Fprintf(&buf, "role:             %s\n", r.node.Role)
	fmt.Fprintf(&buf, "organization:     %s\n", r.node.Organization)
	fmt.Fprintf(&buf, "tokens_staked:    %g\n", r.node.TokensStaked)
	fmt.Fprintf(&buf, "delegators:       %d\n", r.node.Delegators)
	fmt.Fprintf(&buf, "delegators_staked: %g\n", r.node.DelegatorsStaked)
	fmt.Fprintf(&buf, "country:          %s\n", r.node.Country)
	fmt.Fprintf(&buf, "ip_address:       %s\n", r.node.IPAddress)
	fmt.Fprintf(&buf, "epoch:            %d\n", r.node.Epoch)
	return buf.String()
}

func (r *nodeResult) Oneliner() string {
	return fmt.Sprintf("%s %s tokens=%g", r.node.NodeID, r.node.Role, r.node.TokensStaked)
}

func (r *nodeResult) JSON() any { return r.node }

func runGet(args []string, flags *command.GlobalFlags) (command.Result, error) {
	client := command.MustLoadClient()
	resp, err := client.Flow.GetNode().NodeID(args[0]).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("node not found")
	}
	return &nodeResult{node: resp.Data[0]}, nil
}
