package formatter

import (
	"io"
	"text/tabwriter"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type TableFormatter struct{}

func (f *TableFormatter) Format(w io.Writer, data *model.StackInfo) error {
	ew := &errWriter{w: w}

	ew.printf("Stack: %s\n", data.StackName)
	ew.printf("Status: %s\n", data.Status)
	ew.println()

	if ew.err != nil {
		return ew.err
	}

	if len(data.Resources) > 0 {
		ew.printf("Resources (%d)\n", len(data.Resources))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		ew2 := &errWriter{w: tw}
		ew2.printf("LOGICAL ID\tPHYSICAL ID\tTYPE\tSTATUS\n")
		for _, r := range data.Resources {
			ew2.printf("%s\t%s\t%s\t%s\n", r.LogicalID, r.PhysicalID, r.Type, r.Status)
		}
		if ew2.err != nil {
			return ew2.err
		}
		if err := tw.Flush(); err != nil {
			return err
		}
		ew.println()
	}

	if len(data.Outputs) > 0 {
		ew.printf("Outputs (%d)\n", len(data.Outputs))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		ew2 := &errWriter{w: tw}
		ew2.printf("KEY\tVALUE\tEXPORT NAME\n")
		for _, o := range data.Outputs {
			ew2.printf("%s\t%s\t%s\n", o.Key, o.Value, o.ExportName)
		}
		if ew2.err != nil {
			return ew2.err
		}
		if err := tw.Flush(); err != nil {
			return err
		}
		ew.println()
	}

	if len(data.Exports) > 0 {
		ew.printf("Exports (%d)\n", len(data.Exports))
		tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
		ew2 := &errWriter{w: tw}
		ew2.printf("NAME\tVALUE\n")
		for _, e := range data.Exports {
			ew2.printf("%s\t%s\n", e.Name, e.Value)
		}
		if ew2.err != nil {
			return ew2.err
		}
		if err := tw.Flush(); err != nil {
			return err
		}
	}

	return ew.err
}

func (f *TableFormatter) FormatEvents(w io.Writer, data *model.StackEvents) error {
	ew := &errWriter{w: w}
	ew.printf("Events for %s (%d)\n", data.StackName, len(data.Events))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	ew2 := &errWriter{w: tw}
	ew2.printf("TIMESTAMP\tLOGICAL ID\tSTATUS\tREASON\n")
	for _, e := range data.Events {
		ew2.printf("%s\t%s\t%s\t%s\n", e.Timestamp, e.LogicalID, e.Status, e.StatusReason)
	}
	if ew2.err != nil {
		return ew2.err
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	return ew.err
}

func (f *TableFormatter) FormatList(w io.Writer, data *model.StackList) error {
	ew := &errWriter{w: w}
	ew.printf("Stacks (%d)\n", len(data.Stacks))

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	ew2 := &errWriter{w: tw}
	ew2.printf("NAME\tSTATUS\tUPDATED\tDESCRIPTION\n")
	for _, s := range data.Stacks {
		updated := s.UpdatedAt
		if updated == "" {
			updated = s.CreatedAt
		}
		desc := s.Description
		if len(desc) > 60 {
			desc = desc[:57] + "..."
		}
		ew2.printf("%s\t%s\t%s\t%s\n", s.StackName, s.Status, updated, desc)
	}
	if ew2.err != nil {
		return ew2.err
	}
	if err := tw.Flush(); err != nil {
		return err
	}

	return ew.err
}
