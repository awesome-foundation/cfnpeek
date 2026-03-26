package formatter

import (
	"fmt"
	"io"
)

// errWriter wraps an io.Writer and captures the first write error.
// Subsequent writes after an error are no-ops.
type errWriter struct {
	w   io.Writer
	err error
}

func (ew *errWriter) printf(format string, args ...any) {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprintf(ew.w, format, args...)
}

func (ew *errWriter) println() {
	if ew.err != nil {
		return
	}
	_, ew.err = fmt.Fprintln(ew.w)
}
