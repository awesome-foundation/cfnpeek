package formatter

import (
	"fmt"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

// Formatter renders a StackInfo to an output stream.
type Formatter interface {
	Format(w io.Writer, data *model.StackInfo) error
}

// ShortSetter can toggle compact output mode.
// Table, CSV, and INI formatters implement this.
type ShortSetter interface {
	SetShort(short bool)
}

// ListFormatter renders a StackList to an output stream.
type ListFormatter interface {
	FormatList(w io.Writer, data *model.StackList) error
}

// EventFormatter renders a StackEvents to an output stream.
type EventFormatter interface {
	FormatEvents(w io.Writer, data *model.StackEvents) error
}

var registry = map[string]func() Formatter{
	"json":  func() Formatter { return &JSONFormatter{} },
	"yaml":  func() Formatter { return &YAMLFormatter{} },
	"toml":  func() Formatter { return &TOMLFormatter{} },
	"xml":   func() Formatter { return &XMLFormatter{} },
	"ini":   func() Formatter { return &INIFormatter{} },
	"csv":   func() Formatter { return &CSVFormatter{} },
	"table": func() Formatter { return &TableFormatter{} },
}

// Get returns a formatter by name.
func Get(name string) (Formatter, error) {
	ctor, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown format %q (available: %s)", name, Available())
	}
	return ctor(), nil
}

// Available returns a comma-separated list of registered format names.
func Available() string {
	names := []string{"json", "yaml", "toml", "xml", "ini", "csv", "table"}
	result := ""
	for i, n := range names {
		if i > 0 {
			result += ", "
		}
		result += n
	}
	return result
}
