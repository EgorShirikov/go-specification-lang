package printers

import (
	"os"
)

// NewStdoutPrinter creates a printer for stdout
func NewStdoutPrinter() ErrorPrinter {
	return NewStreamPrinter(os.Stdout)
}
