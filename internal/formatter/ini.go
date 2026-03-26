package formatter

import (
	"fmt"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type INIFormatter struct{}

func (f *INIFormatter) Format(w io.Writer, data *model.StackInfo) error {
	fmt.Fprintf(w, "[stack]\n")
	fmt.Fprintf(w, "name = %s\n", data.StackName)
	fmt.Fprintf(w, "id = %s\n", data.StackID)
	fmt.Fprintf(w, "status = %s\n", data.Status)

	if len(data.Resources) > 0 {
		fmt.Fprintf(w, "\n[resources]\n")
		for i, r := range data.Resources {
			fmt.Fprintf(w, "resource.%d.logical_id = %s\n", i, r.LogicalID)
			fmt.Fprintf(w, "resource.%d.physical_id = %s\n", i, r.PhysicalID)
			fmt.Fprintf(w, "resource.%d.type = %s\n", i, r.Type)
			fmt.Fprintf(w, "resource.%d.status = %s\n", i, r.Status)
			fmt.Fprintf(w, "resource.%d.last_updated = %s\n", i, r.LastUpdated)
		}
	}

	if len(data.Outputs) > 0 {
		fmt.Fprintf(w, "\n[outputs]\n")
		for i, o := range data.Outputs {
			fmt.Fprintf(w, "output.%d.key = %s\n", i, o.Key)
			fmt.Fprintf(w, "output.%d.value = %s\n", i, o.Value)
			if o.Description != "" {
				fmt.Fprintf(w, "output.%d.description = %s\n", i, o.Description)
			}
			if o.ExportName != "" {
				fmt.Fprintf(w, "output.%d.export_name = %s\n", i, o.ExportName)
			}
		}
	}

	if len(data.Exports) > 0 {
		fmt.Fprintf(w, "\n[exports]\n")
		for i, e := range data.Exports {
			fmt.Fprintf(w, "export.%d.name = %s\n", i, e.Name)
			fmt.Fprintf(w, "export.%d.value = %s\n", i, e.Value)
		}
	}

	return nil
}

func (f *INIFormatter) FormatList(w io.Writer, data *model.StackList) error {
	for i, s := range data.Stacks {
		if i > 0 {
			fmt.Fprintln(w)
		}
		fmt.Fprintf(w, "[stack.%d]\n", i)
		fmt.Fprintf(w, "name = %s\n", s.StackName)
		fmt.Fprintf(w, "status = %s\n", s.Status)
		fmt.Fprintf(w, "created_at = %s\n", s.CreatedAt)
		if s.UpdatedAt != "" {
			fmt.Fprintf(w, "updated_at = %s\n", s.UpdatedAt)
		}
		if s.Description != "" {
			fmt.Fprintf(w, "description = %s\n", s.Description)
		}
	}
	return nil
}
