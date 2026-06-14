package errors

import (
	"testing"
)

func TestSeverityValues(t *testing.T) {
	tests := []struct {
		name     string
		severity Severity
		expected string
	}{
		{"Critical", SeverityCritical, "crit"},
		{"Warning", SeverityWarning, "warn"},
		{"Info", SeverityInfo, "info"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.severity) != tt.expected {
				t.Errorf("Severity = %v, want %v", tt.severity, tt.expected)
			}
		})
	}
}

func TestPosition(t *testing.T) {
	pos := Position{
		File:   "test.go",
		Line:   42,
		Column: 10,
	}

	if pos.File != "test.go" {
		t.Errorf("Position.File = %v, want test.go", pos.File)
	}
	if pos.Line != 42 {
		t.Errorf("Position.Line = %v, want 42", pos.Line)
	}
	if pos.Column != 10 {
		t.Errorf("Position.Column = %v, want 10", pos.Column)
	}
}

func TestRange(t *testing.T) {
	rng := Range{
		Start: Position{File: "test.go", Line: 1, Column: 1},
		End:   Position{File: "test.go", Line: 5, Column: 20},
	}

	if rng.Start.Line != 1 {
		t.Errorf("Range.Start.Line = %v, want 1", rng.Start.Line)
	}
	if rng.End.Line != 5 {
		t.Errorf("Range.End.Line = %v, want 5", rng.End.Line)
	}
}

func TestError(t *testing.T) {
	err := &Error{
		Position: Position{
			File:   "abs.go",
			Line:   42,
			Column: 10,
		},
		Message:  "Assertion failed: x > 0",
		Severity: SeverityCritical,
	}

	if err.Position.File != "abs.go" {
		t.Errorf("Error.Position.File = %v, want abs.go", err.Position.File)
	}
	if err.Message != "Assertion failed: x > 0" {
		t.Errorf("Error.Message = %v, want 'Assertion failed: x > 0'", err.Message)
	}
	if err.Severity != SeverityCritical {
		t.Errorf("Error.Severity = %v, want %v", err.Severity, SeverityCritical)
	}
}

func TestErrorWithRange(t *testing.T) {
	err := &Error{
		Position: Position{
			File:   "test.go",
			Line:   10,
			Column: 5,
		},
		Range: &Range{
			Start: Position{File: "test.go", Line: 10, Column: 5},
			End:   Position{File: "test.go", Line: 10, Column: 15},
		},
		Message:  "Type mismatch",
		Severity: SeverityWarning,
	}

	if err.Range == nil {
		t.Fatal("Error.Range should not be nil")
	}
	if err.Range.Start.Column != 5 {
		t.Errorf("Error.Range.Start.Column = %v, want 5", err.Range.Start.Column)
	}
	if err.Range.End.Column != 15 {
		t.Errorf("Error.Range.End.Column = %v, want 15", err.Range.End.Column)
	}
}
