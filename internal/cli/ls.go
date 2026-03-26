package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cfnaws "github.com/awesome-foundation/cfnpeek/internal/aws"
	"github.com/awesome-foundation/cfnpeek/internal/formatter"
)

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Short:   "List all active CloudFormation stacks in the region",
		Example: `  cfnpeek ls
  cfnpeek ls -r us-east-1
  cfnpeek ls -f json -p production`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			client, err := newClient(ctx)
			if err != nil {
				return err
			}

			list, err := client.FetchStacks(ctx)
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

			lf, ok := fmtr.(formatter.ListFormatter)
			if !ok {
				return fmt.Errorf("format %q does not support list output", resolved)
			}

			return lf.FormatList(os.Stdout, list)
		},
	}
}
