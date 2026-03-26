package cli

import (
	"context"
	"fmt"
	"os"

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
	// Inspect command flags (local to root, not inherited).
	var (
		showResources bool
		showOutputs   bool
		showExports   bool
		typeFilter    string
		grepFilter    string
	)

	cmd := &cobra.Command{
		Use:   "cfnpeek <stack-name-or-arn> [flags]",
		Short: "Inspect AWS CloudFormation stack resources, outputs, and exports",
		Long: `cfnpeek is a read-only CLI for inspecting deployed AWS CloudFormation stacks.

It displays resources, outputs, and exports in multiple formats.
Output defaults to table when running in a terminal, JSON when piped.

Use "cfnpeek ls" to list all stacks in a region.`,
		Example: `  cfnpeek my-stack
  cfnpeek my-stack --resources
  cfnpeek my-stack --outputs --exports
  cfnpeek my-stack --type ec2
  cfnpeek my-stack --grep vpc
  cfnpeek my-stack --type ec2 --grep vpc
  cfnpeek my-stack -f json -r eu-west-1
  cfnpeek arn:aws:cloudformation:us-east-1:123456789:stack/my-stack/guid`,
		Args:    cobra.ExactArgs(1),
		Version: version,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			stackName := args[0]

			hasTypeFilter := typeFilter != ""
			hasGrepFilter := grepFilter != ""

			// Determine which sections to fetch from AWS.
			// --type implies resources; --grep implies outputs and exports.
			// Explicit section flags are also honoured.
			wantResources := showResources || hasTypeFilter
			wantOutputs := showOutputs || hasGrepFilter
			wantExports := showExports || hasGrepFilter

			// No section flags and no filters = show all.
			if !showResources && !showOutputs && !showExports && !hasTypeFilter && !hasGrepFilter {
				wantResources = true
				wantOutputs = true
				wantExports = true
			}

			client, err := newClient(ctx)
			if err != nil {
				return err
			}

			info, err := client.FetchStackInfo(ctx, stackName, wantResources, wantOutputs, wantExports)
			if err != nil {
				return fmt.Errorf("%s", cfnaws.FormatError(err))
			}

			// Apply post-fetch filters on the model data.
			if hasTypeFilter {
				info.Resources = filter.Resources(info.Resources, typeFilter)
				// When --type is used without --grep, suppress outputs/exports
				// unless they were explicitly requested.
				if !hasGrepFilter && !showOutputs {
					info.Outputs = nil
				}
				if !hasGrepFilter && !showExports {
					info.Exports = nil
				}
			}
			if hasGrepFilter {
				info.Outputs = filter.Outputs(info.Outputs, grepFilter)
				info.Exports = filter.Exports(info.Exports, grepFilter)
				// When --grep is used without --type, suppress resources
				// unless they were explicitly requested.
				if !hasTypeFilter && !showResources {
					info.Resources = nil
				}
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

	// --- Inspect-specific flags (only on root command) ---
	flags := cmd.Flags()
	flags.BoolVar(&showResources, "resources", false, "Show only resources")
	flags.BoolVar(&showOutputs, "outputs", false, "Show only outputs")
	flags.BoolVar(&showExports, "exports", false, "Show only exports")
	flags.StringVar(&typeFilter, "type", "", "Filter resources by type (case-insensitive substring match, implies --resources)")
	flags.StringVar(&grepFilter, "grep", "", "Filter outputs and exports by key/name or value (case-insensitive substring match, implies --outputs --exports)")

	// --- Subcommands ---
	cmd.AddCommand(newLsCmd())

	return cmd
}
