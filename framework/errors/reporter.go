package errors

// ErrorReporter defines the interface for error reporting
type ErrorReporter interface {
	// ReportError reports a single error
	ReportError(err *Error) error

	// ReportErrors reports multiple errors
	ReportErrors(errs []*Error) error

	// SetFormatter changes the formatter dynamically
	SetFormatter(formatter ErrorFormatter)

	// SetPrinter changes the printer dynamically
	SetPrinter(printer ErrorPrinter)

	// Close closes any resources (e.g., file handles)
	Close() error
}

// ErrorPrinter defines the interface for output destinations
type ErrorPrinter interface {
	// Print writes the output without a newline
	Print(output string) error

	// PrintLine writes the output with a newline
	PrintLine(output string) error

	// Close closes any resources
	Close() error
}

// OutputDestination specifies where to send error output
type OutputDestination string

const (
	// DestinationStdout sends output to stdout
	DestinationStdout OutputDestination = "stdout"
	// DestinationStderr sends output to stderr
	DestinationStderr OutputDestination = "stderr"
	// DestinationFile sends output to a file
	DestinationFile OutputDestination = "file"
)

// ReporterConfig configures the ErrorReporter
type ReporterConfig struct {
	// Format is the output format
	Format OutputFormat

	// OutputDestination is where to send output
	OutputDestination OutputDestination

	// FilePath is the path to output file (required if Destination is file)
	FilePath string

	// Append determines whether to append to file or truncate
	Append bool
}
