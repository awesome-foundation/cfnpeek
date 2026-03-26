package formatter

import (
	"io"

	"github.com/BurntSushi/toml"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type TOMLFormatter struct{}

func (f *TOMLFormatter) Format(w io.Writer, data *model.StackInfo) error {
	return toml.NewEncoder(w).Encode(data)
}

func (f *TOMLFormatter) FormatList(w io.Writer, data *model.StackList) error {
	return toml.NewEncoder(w).Encode(data)
}
