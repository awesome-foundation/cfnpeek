package formatter

import (
	"io"

	"gopkg.in/yaml.v3"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type YAMLFormatter struct{}

func (f *YAMLFormatter) Format(w io.Writer, data *model.StackInfo) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return enc.Close()
}

func (f *YAMLFormatter) FormatList(w io.Writer, data *model.StackList) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(data); err != nil {
		return err
	}
	return enc.Close()
}
