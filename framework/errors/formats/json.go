package formats

import (
	"encoding/json"
	"fmt"

	"gpma/framework/errors"
)

// JSONFormatter implements ErrorFormatter for JSON output
type JSONFormatter struct{}

// NewJSONFormatter creates a new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{}
}

// FormatErrors formats a slice of errors as JSON
func (f *JSONFormatter) FormatErrors(errs []*errors.Error, config errors.FormatterConfig) (string, error) {
	if len(errs) == 0 {
		return `{"errors":[]}`, nil
	}

	output := map[string]interface{}{
		"errors": errs,
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal errors to JSON: %w", err)
	}

	return string(data), nil
}

// FormatError formats a single error as JSON
func (f *JSONFormatter) FormatError(err *errors.Error, config errors.FormatterConfig) (string, error) {
	data, marshalErr := json.MarshalIndent(err, "", "  ")
	if marshalErr != nil {
		return "", fmt.Errorf("failed to marshal error to JSON: %w", marshalErr)
	}

	return string(data), nil
}
