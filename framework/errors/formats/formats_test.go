package formats

import (
	"strings"
	"testing"

	"gpma/framework/errors"
)

func createTestError() *errors.Error {
	return &errors.Error{
		Position: errors.Position{
			File:   "abs.go",
			Line:   42,
			Column: 10,
		},
		Message:  "Assertion failed: x > 0",
		Severity: errors.SeverityCritical,
	}
}

func createTestErrorWithRange() *errors.Error {
	return &errors.Error{
		Position: errors.Position{
			File:   "test.go",
			Line:   10,
			Column: 5,
		},
		Range: &errors.Range{
			Start: errors.Position{File: "test.go", Line: 10, Column: 5},
			End:   errors.Position{File: "test.go", Line: 10, Column: 15},
		},
		Message:  "Type mismatch",
		Severity: errors.SeverityWarning,
	}
}

func TestJSONFormatter_FormatError(t *testing.T) {
	formatter := NewJSONFormatter()
	config := errors.FormatterConfig{Format: errors.FormatJSON}

	testErr := createTestError()
	output, formatErr := formatter.FormatError(testErr, config)
	if formatErr != nil {
		t.Fatalf("FormatError failed: %v", formatErr)
	}

	// Check that output contains expected fields
	if !strings.Contains(output, `"file": "abs.go"`) {
		t.Errorf("Output should contain file name")
	}
	if !strings.Contains(output, `"line": 42`) {
		t.Errorf("Output should contain line number")
	}
	if !strings.Contains(output, `"column": 10`) {
		t.Errorf("Output should contain column number")
	}
	if !strings.Contains(output, `"message"`) {
		t.Errorf("Output should contain message field, got: %s", output)
	}
	if !strings.Contains(output, `"severity": "crit"`) {
		t.Errorf("Output should contain severity")
	}
}

func TestJSONFormatter_FormatErrors(t *testing.T) {
	formatter := NewJSONFormatter()
	config := errors.FormatterConfig{Format: errors.FormatJSON}

	errs := []*errors.Error{
		createTestError(),
		createTestErrorWithRange(),
	}

	output, err := formatter.FormatErrors(errs, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	// Check that output contains errors array
	if !strings.Contains(output, `"errors"`) {
		t.Errorf("Output should contain errors array")
	}
	if !strings.Contains(output, `"file": "abs.go"`) {
		t.Errorf("Output should contain first error")
	}
	if !strings.Contains(output, `"file": "test.go"`) {
		t.Errorf("Output should contain second error")
	}
}

func TestJSONFormatter_FormatErrors_Empty(t *testing.T) {
	formatter := NewJSONFormatter()
	config := errors.FormatterConfig{Format: errors.FormatJSON}

	output, err := formatter.FormatErrors([]*errors.Error{}, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	expected := `{"errors":[]}`
	if output != expected {
		t.Errorf("Output = %v, want %v", output, expected)
	}
}

func TestTextFormatter_FormatError(t *testing.T) {
	formatter := NewTextFormatter()
	config := errors.FormatterConfig{Format: errors.FormatText}

	testErr := createTestError()
	output, formatErr := formatter.FormatError(testErr, config)
	if formatErr != nil {
		t.Fatalf("FormatError failed: %v", formatErr)
	}

	expected := "abs.go:42:10: error: Assertion failed: x > 0"
	if output != expected {
		t.Errorf("Output = %v, want %v", output, expected)
	}
}

func TestTextFormatter_FormatErrors(t *testing.T) {
	formatter := NewTextFormatter()
	config := errors.FormatterConfig{Format: errors.FormatText}

	errs := []*errors.Error{
		createTestError(),
		createTestErrorWithRange(),
	}

	output, err := formatter.FormatErrors(errs, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "abs.go:42:10: error: Assertion failed: x > 0") {
		t.Errorf("First line should contain first error")
	}
	if !strings.Contains(lines[1], "test.go:10:5: warning: Type mismatch") {
		t.Errorf("Second line should contain second error")
	}
}

func TestTextFormatter_FormatErrors_Empty(t *testing.T) {
	formatter := NewTextFormatter()
	config := errors.FormatterConfig{Format: errors.FormatText}

	output, err := formatter.FormatErrors([]*errors.Error{}, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	if output != "" {
		t.Errorf("Output should be empty for no errors, got: %v", output)
	}
}

func TestCompactFormatter_FormatError(t *testing.T) {
	formatter := NewCompactFormatter()
	config := errors.FormatterConfig{Format: errors.FormatCompact}

	testErr := createTestError()
	output, formatErr := formatter.FormatError(testErr, config)
	if formatErr != nil {
		t.Fatalf("FormatError failed: %v", formatErr)
	}

	expected := "abs.go:42:10:crit:Assertion failed: x > 0"
	if output != expected {
		t.Errorf("Output = %v, want %v", output, expected)
	}
}

func TestCompactFormatter_FormatErrors(t *testing.T) {
	formatter := NewCompactFormatter()
	config := errors.FormatterConfig{Format: errors.FormatCompact}

	errs := []*errors.Error{
		createTestError(),
		createTestErrorWithRange(),
	}

	output, err := formatter.FormatErrors(errs, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	lines := strings.Split(output, "\n")
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if !strings.Contains(lines[0], "abs.go:42:10:crit:Assertion failed: x > 0") {
		t.Errorf("First line should contain first error")
	}
	if !strings.Contains(lines[1], "test.go:10:5:warn:Type mismatch") {
		t.Errorf("Second line should contain second error")
	}
}

func TestCompactFormatter_FormatErrors_Empty(t *testing.T) {
	formatter := NewCompactFormatter()
	config := errors.FormatterConfig{Format: errors.FormatCompact}

	output, err := formatter.FormatErrors([]*errors.Error{}, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	if output != "" {
		t.Errorf("Output should be empty for no errors, got: %v", output)
	}
}

func TestVSCodeFormatter_FormatError(t *testing.T) {
	formatter := NewVSCodeFormatter()
	config := errors.FormatterConfig{Format: errors.FormatVSCode}

	testErr := createTestError()
	output, formatErr := formatter.FormatError(testErr, config)
	if formatErr != nil {
		t.Fatalf("FormatError failed: %v", formatErr)
	}

	// Check that output contains expected VS Code diagnostic fields
	if !strings.Contains(output, `"message"`) {
		t.Errorf("Output should contain message field")
	}
	if !strings.Contains(output, `"severity"`) {
		t.Errorf("Output should contain severity field")
	}
	if !strings.Contains(output, `"range"`) {
		t.Errorf("Output should contain range field")
	}
	if !strings.Contains(output, `"source"`) {
		t.Errorf("Output should contain source field")
	}
}

func TestVSCodeFormatter_FormatErrors(t *testing.T) {
	formatter := NewVSCodeFormatter()
	config := errors.FormatterConfig{Format: errors.FormatVSCode}

	errs := []*errors.Error{
		createTestError(),
		createTestErrorWithRange(),
	}

	output, err := formatter.FormatErrors(errs, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	// Check that output contains diagnostics array
	if !strings.Contains(output, `"diagnostics"`) {
		t.Errorf("Output should contain diagnostics array")
	}
}

func TestVSCodeFormatter_FormatErrors_Empty(t *testing.T) {
	formatter := NewVSCodeFormatter()
	config := errors.FormatterConfig{Format: errors.FormatVSCode}

	output, err := formatter.FormatErrors([]*errors.Error{}, config)
	if err != nil {
		t.Fatalf("FormatErrors failed: %v", err)
	}

	expected := `{"diagnostics":[]}`
	if output != expected {
		t.Errorf("Output = %v, want %v", output, expected)
	}
}

func TestVSCodeFormatter_SeverityMapping(t *testing.T) {
	formatter := NewVSCodeFormatter()
	config := errors.FormatterConfig{Format: errors.FormatVSCode}

	tests := []struct {
		name     string
		severity errors.Severity
		expected int
	}{
		{"Critical", errors.SeverityCritical, 1},
		{"Warning", errors.SeverityWarning, 2},
		{"Info", errors.SeverityInfo, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testErr := &errors.Error{
				Position: errors.Position{File: "test.go", Line: 1, Column: 1},
				Message:  "Test message",
				Severity: tt.severity,
			}

			output, formatErr := formatter.FormatError(testErr, config)
			if formatErr != nil {
				t.Fatalf("FormatError failed: %v", formatErr)
			}

			// Check that severity is correctly mapped
			expectedSeverity := `"severity": ` + string(rune('0'+tt.expected))
			if !strings.Contains(output, expectedSeverity) {
				t.Errorf("Output should contain severity %d", tt.expected)
			}
		})
	}
}

func TestVSCodeFormatter_RangeHandling(t *testing.T) {
	formatter := NewVSCodeFormatter()
	config := errors.FormatterConfig{Format: errors.FormatVSCode}

	// Test with range
	errWithRange := createTestErrorWithRange()
	output, err := formatter.FormatError(errWithRange, config)
	if err != nil {
		t.Fatalf("FormatError failed: %v", err)
	}

	// VS Code uses 0-based indexing, so line 10 becomes 9, column 5 becomes 4
	if !strings.Contains(output, `"line": 9`) {
		t.Errorf("Output should contain 0-based line number")
	}
	if !strings.Contains(output, `"character": 4`) {
		t.Errorf("Output should contain 0-based column number")
	}

	// Test without range (should default to single character)
	errWithoutRange := createTestError()
	output, err = formatter.FormatError(errWithoutRange, config)
	if err != nil {
		t.Fatalf("FormatError failed: %v", err)
	}

	// Should still have range with single character
	if !strings.Contains(output, `"range"`) {
		t.Errorf("Output should contain range even for single position")
	}
}
