package formatter

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data *model.StackInfo) error {
	fmt.Fprintf(w, "Stack: %s\n", data.StackName)
	fmt.Fprintf(w, "Status: %s\n", data.Status)
	fmt.Fprintln(w)

	if len(data.Resources) > 0 {
		fmt.Fprintf(w, "Resources (%d)\n", len(data.Resources))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "LOGICAL ID\tPHYSICAL ID\tTYPE\tSTATUS\n")
		for _, r := range data.Resources {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.LogicalID, r.PhysicalID, r.Type, r.Status)
		}
		tw.Flush()
		fmt.Fprintln(w)
	}

	if len(data.Outputs) > 0 {
		fmt.Fprintf(w, "Outputs (%d)\n", len(data.Outputs))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "KEY\tVALUE\tEXPORT NAME\n")
		for _, o := range data.Outputs {
			fmt.Fprintf(tw, "%s\t%s\t%s\n", o.Key, o.Value, o.ExportName)
		}
		tw.Flush()
		fmt.Fprintln(w)
	}

	if len(data.Exports) > 0 {
		fmt.Fprintf(w, "Exports (%d)\n", len(data.Exports))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		fmt.Fprintf(tw, "NAME\tVALUE\n")
		for _, e := range data.Exports {
			fmt.Fprintf(tw, "%s\t%s\n", e.Name, e.Value)
		}
		tw.Flush()
	}

	return nil
}

func (f *TableFormatter) FormatList(w io.Writer, data *model.StackList) error {
	fmt.Fprintf(w, "Stacks (%d)\n", len(data.Stacks))
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "NAME\tSTATUS\tUPDATED\tDESCRIPTION\n")
	for _, s := range data.Stacks {
		updated := s.UpdatedAt
		if updated == "" {
			updated = s.CreatedAt
		}
		desc := s.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.StackName, s.Status, updated, desc)
	}
	return tw.Flush()
}
