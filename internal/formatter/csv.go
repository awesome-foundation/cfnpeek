package formatter

import (
	"encoding/csv"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type CSVFormatter struct{}

func (f *CSVFormatter) Format(w io.Writer, data *model.StackInfo) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if len(data.Resources) > 0 {
		if err := cw.Write([]string{"logical_id", "physical_id", "type", "status", "last_updated"}); err != nil {
			return err
		}
		for _, r := range data.Resources {
			if err := cw.Write([]string{r.LogicalID, r.PhysicalID, r.Type, r.Status, r.LastUpdated}); err != nil {
				return err
			}
		}
	}

	if len(data.Outputs) > 0 {
		if len(data.Resources) > 0 {
			cw.Flush() // flush before section separator
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
		}
		if err := cw.Write([]string{"key", "value", "description", "export_name"}); err != nil {
			return err
		}
		for _, o := range data.Outputs {
			if err := cw.Write([]string{o.Key, o.Value, o.Description, o.ExportName}); err != nil {
				return err
			}
		}
	}

	if len(data.Exports) > 0 {
		if len(data.Resources) > 0 || len(data.Outputs) > 0 {
			cw.Flush() // flush before section separator
		if _, err := io.WriteString(w, "\n"); err != nil {
			return err
		}
		}
		if err := cw.Write([]string{"name", "value"}); err != nil {
			return err
		}
		for _, e := range data.Exports {
			if err := cw.Write([]string{e.Name, e.Value}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (f *CSVFormatter) FormatList(w io.Writer, data *model.StackList) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if err := cw.Write([]string{"stack_name", "status", "created_at", "updated_at", "description"}); err != nil {
		return err
	}
	for _, s := range data.Stacks {
		if err := cw.Write([]string{s.StackName, s.Status, s.CreatedAt, s.UpdatedAt, s.Description}); err != nil {
			return err
		}
	}
	return nil
}
