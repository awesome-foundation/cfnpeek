package formatter

import (
	"encoding/csv"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type CSVFormatter struct{ short bool }

func (f *CSVFormatter) SetShort(short bool) { f.short = short }

func (f *CSVFormatter) needsSep(w io.Writer, cw *csv.Writer, prior bool) (bool, error) {
	if prior {
		cw.Flush()
		if _, err := io.WriteString(w, "\n"); err != nil {
			return true, err
		}
	}
	return true, nil
}

func (f *CSVFormatter) Format(w io.Writer, data *model.StackInfo) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	hasPrior := false
	var err error

	if len(data.Resources) > 0 {
		if f.short {
			if err := cw.Write([]string{"logical_id", "type", "status"}); err != nil {
				return err
			}
			for _, r := range data.Resources {
				if err := cw.Write([]string{r.LogicalID, r.Type, r.Status}); err != nil {
					return err
				}
			}
		} else {
			if err := cw.Write([]string{"logical_id", "physical_id", "type", "status", "last_updated"}); err != nil {
				return err
			}
			for _, r := range data.Resources {
				if err := cw.Write([]string{r.LogicalID, r.PhysicalID, r.Type, r.Status, r.LastUpdated}); err != nil {
					return err
				}
			}
		}
		hasPrior = true
	}

	if len(data.Outputs) > 0 {
		if hasPrior, err = f.needsSep(w, cw, hasPrior); err != nil {
			return err
		}
		if f.short {
			if err := cw.Write([]string{"key", "value"}); err != nil {
				return err
			}
			for _, o := range data.Outputs {
				if err := cw.Write([]string{o.Key, o.Value}); err != nil {
					return err
				}
			}
		} else {
			if err := cw.Write([]string{"key", "value", "description", "export_name"}); err != nil {
				return err
			}
			for _, o := range data.Outputs {
				if err := cw.Write([]string{o.Key, o.Value, o.Description, o.ExportName}); err != nil {
					return err
				}
			}
		}
	}

	if len(data.Exports) > 0 {
		if hasPrior, err = f.needsSep(w, cw, hasPrior); err != nil {
			return err
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

	if len(data.Events) > 0 {
		if _, err = f.needsSep(w, cw, hasPrior); err != nil {
			return err
		}
		if f.short {
			if err := cw.Write([]string{"timestamp", "logical_id", "status", "status_reason"}); err != nil {
				return err
			}
			for _, e := range data.Events {
				if err := cw.Write([]string{e.Timestamp, e.LogicalID, e.Status, e.StatusReason}); err != nil {
					return err
				}
			}
		} else {
			if err := cw.Write([]string{"timestamp", "logical_id", "resource_type", "status", "status_reason"}); err != nil {
				return err
			}
			for _, e := range data.Events {
				if err := cw.Write([]string{e.Timestamp, e.LogicalID, e.ResourceType, e.Status, e.StatusReason}); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (f *CSVFormatter) FormatList(w io.Writer, data *model.StackList) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if f.short {
		if err := cw.Write([]string{"stack_name", "status", "updated_at"}); err != nil {
			return err
		}
		for _, s := range data.Stacks {
			updated := s.UpdatedAt
			if updated == "" {
				updated = s.CreatedAt
			}
			if err := cw.Write([]string{s.StackName, s.Status, updated}); err != nil {
				return err
			}
		}
	} else {
		if err := cw.Write([]string{"stack_name", "status", "created_at", "updated_at", "description"}); err != nil {
			return err
		}
		for _, s := range data.Stacks {
			if err := cw.Write([]string{s.StackName, s.Status, s.CreatedAt, s.UpdatedAt, s.Description}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *CSVFormatter) FormatEvents(w io.Writer, data *model.StackEvents) error {
	cw := csv.NewWriter(w)
	defer cw.Flush()

	if f.short {
		if err := cw.Write([]string{"timestamp", "logical_id", "status", "status_reason"}); err != nil {
			return err
		}
		for _, e := range data.Events {
			if err := cw.Write([]string{e.Timestamp, e.LogicalID, e.Status, e.StatusReason}); err != nil {
				return err
			}
		}
	} else {
		if err := cw.Write([]string{"timestamp", "logical_id", "status", "status_reason", "resource_type", "physical_id"}); err != nil {
			return err
		}
		for _, e := range data.Events {
			if err := cw.Write([]string{e.Timestamp, e.LogicalID, e.Status, e.StatusReason, e.ResourceType, e.PhysicalID}); err != nil {
				return err
			}
		}
	}
	return nil
}
