package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cfnaws "github.com/awesome-foundation/cfnpeek/internal/aws"
	"github.com/awesome-foundation/cfnpeek/internal/formatter"
)

func newLogsCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "logs <stack-name-or-arn>",
		Short: "Show CloudFormation stack events (deploy log)",
		Long: `Show CloudFormation stack events for a stack, ordered oldest to newest.

Events are fetched using DescribeStackEvents and displayed in ascending
timestamp order so you can read the deployment log top to bottom.

Use --limit to restrict output to the N most recent events.`,
		Example: `  cfnpeek logs my-stack
  cfnpeek logs my-stack --limit 20
  cfnpeek logs my-stack -f json
  cfnpeek logs arn:aws:cloudformation:us-east-1:123456789:stack/my-stack/guid`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			stackName := args[0]

			client, err := newClient(ctx)
			if err != nil {
				return err
			}

			events, err := client.FetchEvents(ctx, stackName, limit)
			if err != nil {
				return fmt.Errorf("%s", cfnaws.FormatError(err))
			}

			resolved := resolveFormat()
			fmtr, err := formatter.Get(resolved)
			if err != nil {
				return err
			}

			ef, ok := fmtr.(formatter.EventFormatter)
			if !ok {
				return fmt.Errorf("format %q does not support event output", resolved)
			}

			return ef.FormatEvents(os.Stdout, events)
		},
	}

	cmd.Flags().IntVar(&limit, "limit", 0, "Show only the last N events (default: all)")

	return cmd
}
