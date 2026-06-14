# Error System Module

A unified error reporting system for the GPMA framework with support for multiple output formats.

## Overview

This module provides a consistent way to report errors across all plugins with support for different output formats including JSON, text, compact, and VS Code diagnostics.

## Core Types

### Error Structure

```go
type Error struct {
    Position Position      // Location in source code
    Range    *Range        // Optional range in source code
    Message  string        // Human-readable error message
    Severity Severity      // Error severity level
}
```

### Severity Levels

- `SeverityCritical` - Critical errors that prevent further processing
- `SeverityWarning` - Warnings that should be addressed but don't prevent processing
- `SeverityInfo` - Informational messages

### Position and Range

```go
type Position struct {
    File   string  // Path to source file
    Line   int     // 1-based line number
    Column int     // 1-based column number
}

type Range struct {
    Start Position  // Start position
    End   Position  // End position
}
```

## Formatters

The module provides four formatter implementations:

### JSON Formatter

Outputs structured JSON format:

```json
{
  "errors": [
    {
      "position": {
        "file": "abs.go",
        "line": 42,
        "column": 10
      },
      "message": "Assertion failed: x > 0",
      "severity": "crit"
    }
  ]
}
```

### Text Formatter

Human-readable text format:

```
abs.go:42:10: error: Assertion failed: x > 0
```

### Compact Formatter

Single-line compact format:

```
abs.go:42:10:crit: Assertion failed: x > 0
```

### VS Code Formatter

VS Code diagnostics format for IDE integration:

```json
{
  "diagnostics": [
    {
      "range": {
        "start": {"line": 41, "character": 9},
        "end": {"line": 41, "character": 10}
      },
      "severity": 1,
      "message": "Assertion failed: x > 0",
      "source": "gpma"
    }
  ]
}
```

## Usage Example

```go
package main

import (
    "fmt"
    "gpma/framework/errors"
    "gpma/framework/errors/formats"
)

func main() {
    // Create errors
    errs := []*errors.Error{
        {
            Position: errors.Position{
                File:   "abs.go",
                Line:   42,
                Column: 10,
            },
            Message:  "Assertion failed: x > 0",
            Severity: errors.SeverityCritical,
        },
    }

    // Format errors
    formatter := formats.NewJSONFormatter()
    config := errors.FormatterConfig{
        Format: errors.FormatJSON,
    }

    output, err := formatter.FormatErrors(errs, config)
    if err != nil {
        panic(err)
    }

    fmt.Println(output)
}
```

## Design Principles

1. **Independence**: The error system is completely independent with no external dependencies
2. **Simplicity**: Minimal core structures with optional extensions
3. **Flexibility**: Plugins can specify either a single position or a range
4. **Extensibility**: Easy to add new formatters without modifying core types

## Testing

Run tests with:

```bash
cd framework/errors
go test ./...
```

## File Structure

```
framework/errors/
├── errors.go              # Core error types
├── formatter.go           # Formatter interface and config
├── errors_test.go         # Tests for error types
├── formats/
│   ├── json.go            # JSON formatter
│   ├── text.go            # Text formatter
│   ├── compact.go         # Compact formatter
│   ├── vscode.go          # VS Code formatter
│   └── formats_test.go    # Tests for formatters
└── README.md              # This file
```
