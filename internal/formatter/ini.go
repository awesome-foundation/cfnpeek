package formatter

import (
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type INIFormatter struct{}

func (f *INIFormatter) Format(w io.Writer, data *model.StackInfo) error {
	ew := &errWriter{w: w}

	ew.printf("[stack]\n")
	ew.printf("name = %s\n", data.StackName)
	ew.printf("id = %s\n", data.StackID)
	ew.printf("status = %s\n", data.Status)

	if len(data.Resources) > 0 {
		ew.printf("\n[resources]\n")
		for i, r := range data.Resources {
			ew.printf("resource.%d.logical_id = %s\n", i, r.LogicalID)
			ew.printf("resource.%d.physical_id = %s\n", i, r.PhysicalID)
			ew.printf("resource.%d.type = %s\n", i, r.Type)
			ew.printf("resource.%d.status = %s\n", i, r.Status)
			ew.printf("resource.%d.last_updated = %s\n", i, r.LastUpdated)
		}
	}

	if len(data.Outputs) > 0 {
		ew.printf("\n[outputs]\n")
		for i, o := range data.Outputs {
			ew.printf("output.%d.key = %s\n", i, o.Key)
			ew.printf("output.%d.value = %s\n", i, o.Value)
			if o.Description != "" {
				ew.printf("output.%d.description = %s\n", i, o.Description)
			}
			if o.ExportName != "" {
				ew.printf("output.%d.export_name = %s\n", i, o.ExportName)
			}
		}
	}

	if len(data.Exports) > 0 {
		ew.printf("\n[exports]\n")
		for i, e := range data.Exports {
			ew.printf("export.%d.name = %s\n", i, e.Name)
			ew.printf("export.%d.value = %s\n", i, e.Value)
		}
	}

	return ew.err
}

func (f *INIFormatter) FormatList(w io.Writer, data *model.StackList) error {
	ew := &errWriter{w: w}

	for i, s := range data.Stacks {
		if i > 0 {
			ew.println()
		}
		ew.printf("[stack.%d]\n", i)
		ew.printf("name = %s\n", s.StackName)
		ew.printf("status = %s\n", s.Status)
		ew.printf("created_at = %s\n", s.CreatedAt)
		if s.UpdatedAt != "" {
			ew.printf("updated_at = %s\n", s.UpdatedAt)
		}
		if s.Description != "" {
			ew.printf("description = %s\n", s.Description)
		}
	}

	return ew.err
}
