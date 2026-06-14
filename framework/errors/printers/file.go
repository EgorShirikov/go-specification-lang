package printers

import (
	"os"
)

// NewFilePrinter creates a printer for a file
func NewFilePrinter(path string, append bool) (ErrorPrinter, error) {
	flags := os.O_CREATE | os.O_WRONLY
	if append {
		flags |= os.O_APPEND
	} else {
		flags |= os.O_TRUNC
	}

	file, err := os.OpenFile(path, flags, 0644)
	if err != nil {
		return nil, err
	}

	return NewStreamPrinterWithCloser(file), nil
}
