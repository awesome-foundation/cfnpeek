package formatter

import (
	"io"

	"github.com/BurntSushi/toml"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type TOMLFormatter struct{}

func (f *TOMLFormatter) Format(w io.Writer, data *model.StackInfo) error {
	doc := map[string]any{
		"stack": map[string]any{
			"name":   data.StackName,
			"id":     data.StackID,
			"status": data.Status,
		},
	}

	if len(data.Resources) > 0 {
		resources := make(map[string]any, len(data.Resources))
		for _, r := range data.Resources {
			resources[r.LogicalID] = map[string]any{
				"physical_id":  r.PhysicalID,
				"type":         r.Type,
				"status":       r.Status,
				"last_updated": r.LastUpdated,
			}
		}
		doc["resources"] = resources
	}

	if len(data.Outputs) > 0 {
		outputs := make(map[string]any, len(data.Outputs))
		for _, o := range data.Outputs {
			outputs[o.Key] = o.Value
		}
		doc["outputs"] = outputs
	}

	if len(data.Exports) > 0 {
		exports := make(map[string]any, len(data.Exports))
		for _, e := range data.Exports {
			exports[e.Name] = e.Value
		}
		doc["exports"] = exports
	}

	return toml.NewEncoder(w).Encode(doc)
}

func (f *TOMLFormatter) FormatList(w io.Writer, data *model.StackList) error {
	doc := make(map[string]any, len(data.Stacks))
	for _, s := range data.Stacks {
		entry := map[string]any{
			"status":     s.Status,
			"created_at": s.CreatedAt,
		}
		if s.UpdatedAt != "" {
			entry["updated_at"] = s.UpdatedAt
		}
		if s.Description != "" {
			entry["description"] = s.Description
		}
		doc[s.StackName] = entry
	}
	return toml.NewEncoder(w).Encode(doc)
}
