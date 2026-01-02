// Package errorsx provides custom error types for compilation,
// lexer, and parser phases of the generator.
package errorsx

import (
	"fmt"
	"os"
)

// Phase represents the stage of compilation where an error occurred.
type Phase string

const (
	PhaseCompilation Phase = "Compilation"
	PhaseLexer       Phase = "Lexer"
	PhaseParser      Phase = "Parser"
)

// Error represents a structured generator error.
type Error struct {
	Phase   Phase  // Lexer / Parser / Compilation
	Message string // Human-readable error message
	Line    int    // Source line number (if available)
	Column  int    // Source column number (if available)
}

// Error implements the built-in error interface.
func (e *Error) Error() string {
	loc := ""
	if e.Line > 0 {
		loc = fmt.Sprintf(" (line %d, col %d)", e.Line, e.Column)
	}
	return fmt.Sprintf("[%s Error]%s: %s", e.Phase, loc, e.Message)
}

// NewCompilationError returns a compilation error.
func NewCompilationError(msg string, l ...int) *Error {
	var line, col int
	if len(l) > 0 {
		line = l[0]
	}
	if len(l) > 1 {
		col = l[1]
	}
	return &Error{Phase: PhaseCompilation, Message: msg, Line: line, Column: col}
}

// NewLexerError returns a lexer error with position.
func NewLexerError(msg string, l ...int) *Error {
	var line, col int
	if len(l) > 0 {
		line = l[0]
	}
	if len(l) > 1 {
		col = l[1]
	}
	return &Error{Phase: PhaseLexer, Message: msg, Line: line, Column: col}
}

// NewParserError returns a parser error with position.
func NewParserError(msg string, l ...int) *Error {
	var line, col int
	if len(l) > 0 {
		line = l[0]
	}
	if len(l) > 1 {
		col = l[1]
	}
	return &Error{Phase: PhaseParser, Message: msg, Line: line, Column: col}
}

// PanicLexerError panics with a lexer error.
func PanicCompilationError(msg string, path string, line int, col int) {
	printSourceContext(
		path,
		line,
		col,
		msg,
		PhaseCompilation,
	)
	os.Exit(1)
}

// PanicLexerError panics with a lexer error.
func PanicLexerError(msg string, path string, line int, col int) {
	printSourceContext(
		path,
		line,
		col,
		msg,
		PhaseLexer,
	)
	os.Exit(1)
}

// PanicParserError panics with a parser error.
func PanicParserError(msg string, path string, line int, col int) {
	printSourceContext(
		path,
		line,
		col,
		msg,
		PhaseParser,
	)
	os.Exit(1)
}
