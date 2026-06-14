package printers

import (
	"os"
)

// NewStderrPrinter creates a printer for stderr
func NewStderrPrinter() ErrorPrinter {
	return NewStreamPrinter(os.Stderr)
}
