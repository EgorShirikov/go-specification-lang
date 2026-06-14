package errors

// ErrorReporterImpl implements the ErrorReporter interface
type ErrorReporterImpl struct {
	formatter ErrorFormatter
	printer   ErrorPrinter
	config    ReporterConfig
}

// NewErrorReporter creates a new ErrorReporter with the given configuration
func NewErrorReporter(config ReporterConfig) (*ErrorReporterImpl, error) {
	formatter, err := getFormatter(config.Format)
	if err != nil {
		return nil, err
	}
	printer, err := getPrinter(config.OutputDestination, config.FilePath, config.Append)
	if err != nil {
		return nil, err
	}

	return &ErrorReporterImpl{
		formatter: formatter,
		printer:   printer,
		config:    config,
	}, nil
}

// NewErrorReporterWithFormatterPrinter creates a new ErrorReporter with custom formatter and printer
func NewErrorReporterWithFormatterPrinter(formatter ErrorFormatter, printer ErrorPrinter) *ErrorReporterImpl {
	return &ErrorReporterImpl{
		formatter: formatter,
		printer:   printer,
		config: ReporterConfig{
			Format: FormatText, // Default format
		},
	}
}

// ReportError reports a single error
func (r *ErrorReporterImpl) ReportError(err *Error) error {
	formatted, formatErr := r.formatter.FormatError(err, FormatterConfig{Format: r.config.Format})
	if formatErr != nil {
		return formatErr
	}
	return r.printer.PrintLine(formatted)
}

// ReportErrors reports multiple errors
func (r *ErrorReporterImpl) ReportErrors(errs []*Error) error {
	formatted, formatErr := r.formatter.FormatErrors(errs, FormatterConfig{Format: r.config.Format})
	if formatErr != nil {
		return formatErr
	}
	return r.printer.PrintLine(formatted)
}

// SetFormatter changes the formatter dynamically
func (r *ErrorReporterImpl) SetFormatter(formatter ErrorFormatter) {
	r.formatter = formatter
}

// SetPrinter changes the printer dynamically
func (r *ErrorReporterImpl) SetPrinter(printer ErrorPrinter) {
	r.printer = printer
}

// Close closes any resources (e.g., file handles)
func (r *ErrorReporterImpl) Close() error {
	return r.printer.Close()
}

// SetFormatterFactory sets the formatter factory function
// This should be called from the formatter_factory package's init function
func SetFormatterFactory(factory func(OutputFormat) (ErrorFormatter, error)) {
	formatterFactory = factory
}

var formatterFactory func(OutputFormat) (ErrorFormatter, error)

// getFormatter returns the appropriate formatter for the given format
func getFormatter(format OutputFormat) (ErrorFormatter, error) {
	if formatterFactory != nil {
		return formatterFactory(format)
	}
	// Fallback to basic text formatter
	return &basicTextFormatter{}, nil
}

// basicTextFormatter is a simple fallback formatter
type basicTextFormatter struct{}

func (f *basicTextFormatter) FormatError(err *Error, config FormatterConfig) (string, error) {
	return "[" + err.Source + "] " + err.Position.File + ":" + string(rune(err.Position.Line)) + ":" + string(rune(err.Position.Column)) + ": " + err.Message, nil
}

func (f *basicTextFormatter) FormatErrors(errs []*Error, config FormatterConfig) (string, error) {
	var result string
	for _, err := range errs {
		formatted, _ := f.FormatError(err, config)
		result += formatted + "\n"
	}
	return result, nil
}
