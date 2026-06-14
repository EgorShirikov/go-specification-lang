package errors

// Severity represents the severity level of an error
type Severity string

const (
	// SeverityCritical represents a critical error that prevents further processing
	SeverityCritical Severity = "crit"
	// SeverityWarning represents a warning that should be addressed but doesn't prevent processing
	SeverityWarning Severity = "warn"
	// SeverityInfo represents an informational message
	SeverityInfo Severity = "info"
)

// Position represents a position in source code using line:column
type Position struct {
	// File is the path to the source file
	File string `json:"file"`

	// Line is the 1-based line number
	Line int `json:"line"`

	// Column is the 1-based column number
	Column int `json:"column"`
}

// Range represents a range in source code (optional)
type Range struct {
	// Start is the start position
	Start Position `json:"start"`

	// End is the end position
	End Position `json:"end"`
}

// Error represents a unified error structure
type Error struct {
	// Position is the location in source code (single point)
	Position Position `json:"position"`

	// Range is the optional range in source code (if specified by plugin)
	Range *Range `json:"range,omitempty"`

	// Message is the human-readable error message
	Message string `json:"message"`

	// Severity is the error severity level
	Severity Severity `json:"severity"`

	// Source identifies which plugin generated this error (must be set by plugin)
	Source string `json:"source"`
}
