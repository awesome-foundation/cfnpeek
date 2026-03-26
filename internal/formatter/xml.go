package formatter

import (
	"encoding/xml"
	"io"

	"github.com/awesome-foundation/cfnpeek/internal/model"
)

type XMLFormatter struct{}

func (f *XMLFormatter) Format(w io.Writer, data *model.StackInfo) error {
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	defer enc.Close()
	return enc.Encode(data)
}

func (f *XMLFormatter) FormatList(w io.Writer, data *model.StackList) error {
	if _, err := io.WriteString(w, xml.Header); err != nil {
		return err
	}
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	defer enc.Close()
	return enc.Encode(data)
}
