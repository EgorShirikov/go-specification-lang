package formatter_factory

import (
	"fmt"
	"gpma/framework/errors"
	"gpma/framework/errors/formats"
)

// FormatterFactory creates formatter instances without creating import cycles
type FormatterFactory struct{}

// NewFormatterFactory creates a new formatter factory
func NewFormatterFactory() *FormatterFactory {
	return &FormatterFactory{}
}

// GetFormatter returns a formatter for the given format
func (f *FormatterFactory) GetFormatter(format errors.OutputFormat) (errors.ErrorFormatter, error) {
	switch format {
	case errors.FormatJSON:
		return &FormatterAdapter{formatter: formats.NewJSONFormatter()}, nil
	case errors.FormatText:
		return &FormatterAdapter{formatter: formats.NewTextFormatter()}, nil
	case errors.FormatCompact:
		return &FormatterAdapter{formatter: formats.NewCompactFormatter()}, nil
	case errors.FormatVSCode:
		return &FormatterAdapter{formatter: formats.NewVSCodeFormatter()}, nil
	default:
		return nil, &errors.ReporterError{
			Message: fmt.Sprintf("unknown error format: %q (valid formats: text, json, compact, vscode)", format),
		}
	}
}

// FormatterAdapter wraps existing formatters to implement errors.ErrorFormatter
type FormatterAdapter struct {
	formatter interface{}
}

// FormatErrors formats a slice of errors
func (a *FormatterAdapter) FormatErrors(errs []*errors.Error, config errors.FormatterConfig) (string, error) {
	switch f := a.formatter.(type) {
	case interface {
		FormatErrors([]*errors.Error, errors.FormatterConfig) (string, error)
	}:
		return f.FormatErrors(errs, config)
	default:
		return "", &errors.ReporterError{Message: "formatter does not support FormatErrors"}
	}
}

// FormatError formats a single error
func (a *FormatterAdapter) FormatError(err *errors.Error, config errors.FormatterConfig) (string, error) {
	switch f := a.formatter.(type) {
	case interface {
		FormatError(*errors.Error, errors.FormatterConfig) (string, error)
	}:
		return f.FormatError(err, config)
	default:
		return "", &errors.ReporterError{Message: "formatter does not support FormatError"}
	}
}

// GetFormatter is a package-level function that can be called from errors package
func GetFormatter(format errors.OutputFormat) (errors.ErrorFormatter, error) {
	factory := NewFormatterFactory()
	return factory.GetFormatter(format)
}

func init() {
	// Set the default formatter factory for the errors package
	errors.SetFormatterFactory(GetFormatter)
}
