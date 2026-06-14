package formats

import (
	"fmt"
	"strings"

	"gpma/framework/errors"
)

// CompactFormatter implements ErrorFormatter for compact single-line output
type CompactFormatter struct{}

// NewCompactFormatter creates a new compact formatter
func NewCompactFormatter() *CompactFormatter {
	return &CompactFormatter{}
}

// FormatErrors formats a slice of errors in compact format
func (f *CompactFormatter) FormatErrors(errs []*errors.Error, config errors.FormatterConfig) (string, error) {
	if len(errs) == 0 {
		return "", nil
	}

	var builder strings.Builder
	for _, err := range errs {
		line, err := f.FormatError(err, config)
		if err != nil {
			return "", err
		}
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return strings.TrimSpace(builder.String()), nil
}

// FormatError formats a single error in compact format
func (f *CompactFormatter) FormatError(err *errors.Error, config errors.FormatterConfig) (string, error) {
	position := fmt.Sprintf("%s:%d:%d", err.Position.File, err.Position.Line, err.Position.Column)

	return fmt.Sprintf("%s:%s:%s", position, err.Severity, err.Message), nil
}
