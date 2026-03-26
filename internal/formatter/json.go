package formatter

import (
	"encoding/json"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(w io.Writer, data *model.StackInfo) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func (f *JSONFormatter) FormatList(w io.Writer, data *model.StackList) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
