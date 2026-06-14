package printers

import (
	"io"
)

// ErrorPrinter defines the interface for output destinations
type ErrorPrinter interface {
	// Print writes the output without a newline
	Print(output string) error

	// PrintLine writes the output with a newline
	PrintLine(output string) error

	// Close closes any resources
	Close() error
}

// StreamPrinter writes to any io.Writer (stdout, stderr, file, etc.)
type StreamPrinter struct {
	writer io.Writer
	closer io.Closer
}

// NewStreamPrinter creates a printer that writes to any io.Writer
// If the writer is also a Closer, it will be closed when Close() is called
func NewStreamPrinter(writer io.Writer) ErrorPrinter {
	return &StreamPrinter{
		writer: writer,
		closer: nil,
	}
}

// NewStreamPrinterWithCloser creates a printer with a closer (for files, etc.)
func NewStreamPrinterWithCloser(writer io.WriteCloser) ErrorPrinter {
	return &StreamPrinter{
		writer: writer,
		closer: writer,
	}
}

// Print writes the output without a newline
func (p *StreamPrinter) Print(output string) error {
	_, err := p.writer.Write([]byte(output))
	return err
}

// PrintLine writes the output with a newline
func (p *StreamPrinter) PrintLine(output string) error {
	_, err := p.writer.Write([]byte(output + "\n"))
	return err
}

// Close closes any resources
func (p *StreamPrinter) Close() error {
	if p.closer != nil {
		return p.closer.Close()
	}
	return nil
}
