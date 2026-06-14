package errors

import (
	"io"

	"gpma/framework/errors/printers"
)

// getPrinter creates an ErrorPrinter based on the destination
func getPrinter(dest OutputDestination, filePath string, append bool) (ErrorPrinter, error) {
	switch dest {
	case DestinationStdout:
		return printers.NewStdoutPrinter(), nil
	case DestinationStderr:
		return printers.NewStderrPrinter(), nil
	case DestinationFile:
		if filePath == "" {
			return nil, &ReporterError{Message: "file path is required for file output destination"}
		}
		return printers.NewFilePrinter(filePath, append)
	default:
		return nil, &ReporterError{Message: "unknown output destination: " + string(dest)}
	}
}

// NewStreamPrinter creates a printer that writes to any io.Writer
func NewStreamPrinter(writer io.Writer) ErrorPrinter {
	return printers.NewStreamPrinter(writer)
}

// NewStreamPrinterWithCloser creates a printer with a closer (for files, etc.)
func NewStreamPrinterWithCloser(writer io.WriteCloser) ErrorPrinter {
	return printers.NewStreamPrinterWithCloser(writer)
}

// NewStdoutPrinter creates a printer for stdout
func NewStdoutPrinter() ErrorPrinter {
	return printers.NewStdoutPrinter()
}

// NewStderrPrinter creates a printer for stderr
func NewStderrPrinter() ErrorPrinter {
	return printers.NewStderrPrinter()
}

// NewFilePrinter creates a printer for a file
func NewFilePrinter(path string, append bool) (ErrorPrinter, error) {
	return printers.NewFilePrinter(path, append)
}

// ReporterError represents an error in the reporter system
type ReporterError struct {
	Message string
}

func (e *ReporterError) Error() string {
	return e.Message
}

// IsReporterError checks if an error is a ReporterError
func IsReporterError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*ReporterError)
	return ok
}
