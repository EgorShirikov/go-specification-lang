package errors

// OutputFormat specifies the output format
type OutputFormat string

const (
	// FormatJSON represents JSON output format
	FormatJSON OutputFormat = "json"
	// FormatText represents human-readable text output format
	FormatText OutputFormat = "text"
	// FormatCompact represents compact single-line output format
	FormatCompact OutputFormat = "compact"
	// FormatVSCode represents VS Code diagnostics format
	FormatVSCode OutputFormat = "vscode"
)

// FormatterConfig configures error formatting
type FormatterConfig struct {
	// Format is the output format
	Format OutputFormat
}

// ErrorFormatter formats errors for output
type ErrorFormatter interface {
	// FormatErrors formats a slice of errors according to config
	FormatErrors(errors []*Error, config FormatterConfig) (string, error)

	// FormatError formats a single error
	FormatError(err *Error, config FormatterConfig) (string, error)
}
