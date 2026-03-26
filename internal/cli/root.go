package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	cfnaws "github.com/awesome-foundation/cfnpeek/internal/aws"
	"github.com/awesome-foundation/cfnpeek/internal/filter"
	"github.com/awesome-foundation/cfnpeek/internal/formatter"
)

// Global flags (shared by all commands via PersistentFlags).
var (
	region  string
	profile string
	format  string
)

var validSections = map[string]bool{
	"all":       true,
	"resources": true,
	"outputs":   true,
	"exports":   true,
	"events":    true,
	"logs":      true, // silent alias for events
}

// resolveFormat picks the output format, auto-detecting TTY when set to "auto".
func resolveFormat() string {
	if format != "auto" {
		return format
	}
	if isatty.IsTerminal(os.Stdout.Fd()) || isatty.IsCygwinTerminal(os.Stdout.Fd()) {
		return "table"
	}
	return "json"
}

// newClient is a shared helper that all subcommands use.
func newClient(ctx context.Context) (*cfnaws.Client, error) {
	client, err := cfnaws.NewClient(ctx, region, profile)
	if err != nil {
		return nil, fmt.Errorf("%s", cfnaws.FormatError(err))
	}
	return client, nil
}

func NewRootCmd(version string) *cobra.Command {
	var (
		typeFilter string
		grepFilter string
	)

	cmd := &cobra.Command{
		Use:   "cfnpeek <stack> [command] [flags]",
		Short: "Inspect AWS CloudFormation stack resources, outputs, and exports",
		Long: `cfnpeek is a read-only CLI for inspecting deployed AWS CloudFormation stacks.

It displays resources, outputs, and exports in multiple formats.
Output defaults to table when running in a terminal, JSON when piped.

Commands:
  cfnpeek ls                     List all stacks in the region
  cfnpeek <stack>                Show all sections (default)
  cfnpeek <stack> resources      Show resources
  cfnpeek <stack> outputs        Show outputs
  cfnpeek <stack> exports        Show exports
  cfnpeek <stack> events         Show stack events (deploy log)
  cfnpeek <stack> all            Show all sections (explicit)`,
		Example: `  cfnpeek my-stack
  cfnpeek my-stack resources
  cfnpeek my-stack outputs
  cfnpeek my-stack events
  cfnpeek my-stack resources,events
  cfnpeek my-stack events --limit 10
  cfnpeek my-stack resources --type ec2
  cfnpeek my-stack outputs --grep vpc
  cfnpeek my-stack -f json -r eu-west-1
  cfnpeek ls`,
		Args:    cobra.RangeArgs(1, 2),
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			stackName := args[0]
			sectionArg := "all"
			if len(args) > 1 {
				sectionArg = strings.ToLower(args[1])
			}

			// Parse comma-separated sections
			parts := strings.Split(sectionArg, ",")
			want := make(map[string]bool)
			for _, p := range parts {
				p = strings.TrimSpace(p)
				if p == "logs" {
					p = "events" // silent alias
				}
				if !validSections[p] {
					return fmt.Errorf("unknown command %q (available: all, resources, outputs, exports, events)", p)
				}
				if p == "all" {
					want["resources"] = true
					want["outputs"] = true
					want["exports"] = true
					want["events"] = true
				} else {
					want[p] = true
				}
			}

			ctx := context.Background()

			client, err := newClient(ctx)
			if err != nil {
				return err
			}

			// --type and --grep expand what we fetch
			if typeFilter != "" {
				want["resources"] = true
			}
			if grepFilter != "" {
				want["outputs"] = true
				want["exports"] = true
			}

			// Events-only mode uses the standalone events formatter
			if want["events"] && !want["resources"] && !want["outputs"] && !want["exports"] {
				return runEvents(ctx, client, cmd, stackName)
			}

			// Fetch stack info (always needed if we get here)
			info, err := client.FetchStackInfo(ctx, stackName, want["resources"], want["outputs"], want["exports"])
			if err != nil {
				return fmt.Errorf("%s", cfnaws.FormatError(err))
			}

			if typeFilter != "" {
				info.Resources = filter.Resources(info.Resources, typeFilter)
			}
			if grepFilter != "" {
				info.Outputs = filter.Outputs(info.Outputs, grepFilter)
				info.Exports = filter.Exports(info.Exports, grepFilter)
			}

			// Fetch events and embed in StackInfo for unified output
			if want["events"] {
				events, err := client.FetchEvents(ctx, stackName, eventsLimit)
				if err != nil {
					return fmt.Errorf("%s", cfnaws.FormatError(err))
				}
				info.Events = events.Events
			}

			resolved := resolveFormat()
			fmtr, err := formatter.Get(resolved)
			if err != nil {
				return err
			}

			return fmtr.Format(os.Stdout, info)
		},
	}

	// --- Global flags (inherited by all subcommands) ---
	pflags := cmd.PersistentFlags()
	pflags.StringVarP(&region, "region", "r", "", "AWS region (overrides AWS_REGION / config)")
	pflags.StringVarP(&profile, "profile", "p", "", "AWS profile (overrides AWS_PROFILE)")
	pflags.StringVarP(&format, "format", "f", "auto", "Output format: auto, table, json, yaml, toml, xml, ini, csv")

	// --- Inspect flags ---
	flags := cmd.Flags()
	flags.StringVar(&typeFilter, "type", "", "Filter resources by type (case-insensitive substring match)")
	flags.StringVar(&grepFilter, "grep", "", "Filter outputs/exports by key/name or value (case-insensitive substring match)")

	// --- Events flags ---
	addEventsFlags(cmd)

	// --- Subcommands ---
	cmd.AddCommand(newLsCmd())

	return cmd
}
