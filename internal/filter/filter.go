// Package filter provides post-fetch filtering of StackInfo model data.
package filter

import (
	"strings"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

// Resources returns a copy of resources whose Type contains pattern
// (case-insensitive substring match).
func Resources(resources []model.Resource, pattern string) []model.Resource {
	lower := strings.ToLower(pattern)
	var out []model.Resource
	for _, r := range resources {
		if strings.Contains(strings.ToLower(r.Type), lower) {
			out = append(out, r)
		}
	}
	return out
}

// Outputs returns a copy of outputs whose Key or Value contains pattern
// (case-insensitive substring match).
func Outputs(outputs []model.Output, pattern string) []model.Output {
	lower := strings.ToLower(pattern)
	var out []model.Output
	for _, o := range outputs {
		if strings.Contains(strings.ToLower(o.Key), lower) ||
			strings.Contains(strings.ToLower(o.Value), lower) {
			out = append(out, o)
		}
	}
	return out
}

// Exports returns a copy of exports whose Name or Value contains pattern
// (case-insensitive substring match).
func Exports(exports []model.Export, pattern string) []model.Export {
	lower := strings.ToLower(pattern)
	var out []model.Export
	for _, e := range exports {
		if strings.Contains(strings.ToLower(e.Name), lower) ||
			strings.Contains(strings.ToLower(e.Value), lower) {
			out = append(out, e)
		}
	}
	return out
}
