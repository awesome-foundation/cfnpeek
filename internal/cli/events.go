package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cfnaws "github.com/awesome-foundation/cfnpeek/internal/aws"
	"github.com/awesome-foundation/cfnpeek/internal/formatter"
)

var eventsLimit int

func addEventsFlags(cmd *cobra.Command) {
	cmd.Flags().IntVar(&eventsLimit, "limit", 20, "Max number of events to show (0 for all)")
}

func runEvents(ctx context.Context, client *cfnaws.Client, cmd *cobra.Command, stackName string) error {
	limit := eventsLimit

	events, err := client.FetchEvents(ctx, stackName, limit)
	if err != nil {
		return fmt.Errorf("%s", cfnaws.FormatError(err))
	}

	resolved := resolveFormat()
	fmtr, err := formatter.Get(resolved)
	if err != nil {
		return err
	}
	if ss, ok := fmtr.(formatter.ShortSetter); ok {
		ss.SetShort(short)
	}

	ef, ok := fmtr.(formatter.EventFormatter)
	if !ok {
		return fmt.Errorf("format %q does not support event output", resolved)
	}

	return ef.FormatEvents(os.Stdout, events)
}
