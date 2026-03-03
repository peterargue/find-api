package blocks

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/tabwriter"

	"github.com/peterargue/find-api/cmd/findapi/internal/command"
	"github.com/peterargue/find-api/flow"
	"github.com/spf13/cobra"
)

type serviceEventsFlags struct {
	Limit  int `flag:"limit"  info:"Number of events to return (max 100)"`
	Offset int `flag:"offset" info:"Pagination offset"`
}

var serviceEventsFlagsVal = &serviceEventsFlags{}

var serviceEventsCmd = &command.Command{
	Cmd: &cobra.Command{
		Use:     "service-events <height>",
		Short:   "List service events for a block",
		Example: "find blocks service-events 12345678",
		Args:    cobra.ExactArgs(1),
	},
	Flags: serviceEventsFlagsVal,
	Run:   runServiceEvents,
}

type serviceEventsResult struct {
	events []flow.BlockServiceEvent
}

func (r *serviceEventsResult) String() string {
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "HEIGHT\tNAME\tTIMESTAMP")
	for _, e := range r.events {
		fmt.Fprintf(w, "%d\t%s\t%s\n", e.BlockHeight, e.Name, e.Timestamp)
	}
	w.Flush()
	return buf.String()
}

func (r *serviceEventsResult) Oneliner() string {
	return fmt.Sprintf("%d events", len(r.events))
}

func (r *serviceEventsResult) JSON() any { return r.events }

func runServiceEvents(args []string, flags *command.GlobalFlags) (command.Result, error) {
	height, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid height %q: %w", args[0], err)
	}
	client := command.MustLoadClient()
	b := client.Flow.GetBlockServiceEvents().Height(height)
	if serviceEventsFlagsVal.Limit > 0 {
		b = b.Limit(serviceEventsFlagsVal.Limit)
	}
	if serviceEventsFlagsVal.Offset > 0 {
		b = b.Offset(serviceEventsFlagsVal.Offset)
	}
	resp, err := b.Do(context.Background())
	if err != nil {
		return nil, err
	}
	return &serviceEventsResult{events: resp.Data}, nil
}
