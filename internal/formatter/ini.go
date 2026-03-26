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

	for _, r := range data.Resources {
		ew.printf("\n[resource.%s]\n", r.LogicalID)
		ew.printf("physical_id = %s\n", r.PhysicalID)
		ew.printf("type = %s\n", r.Type)
		ew.printf("status = %s\n", r.Status)
		ew.printf("last_updated = %s\n", r.LastUpdated)
	}

	if len(data.Outputs) > 0 {
		ew.printf("\n[outputs]\n")
		for _, o := range data.Outputs {
			ew.printf("%s = %s\n", o.Key, o.Value)
		}
	}

	if len(data.Exports) > 0 {
		ew.printf("\n[exports]\n")
		for _, e := range data.Exports {
			ew.printf("%s = %s\n", e.Name, e.Value)
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
		ew.printf("[%s]\n", s.StackName)
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
