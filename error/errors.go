// Package errorsx provides custom error types for compilation,
// lexer, and parser phases of the compiler.
package errorsx

import (
	"fmt"
)

// Phase represents the stage of compilation where an error occurred.
type Phase string

const (
	PhaseCompilation Phase = "Compilation"
	PhaseLexer       Phase = "Lexer"
	PhaseParser      Phase = "Parser"
)

// Error represents a structured compiler error.
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

// PanicCompilationError panics with a compilation error.
func PanicCompilationError(msg string, l ...int) {
	panic(NewCompilationError(msg, l...))
}

// PanicLexerError panics with a lexer error.
func PanicLexerError(msg string, l ...int) {
	panic(NewLexerError(msg, l...))
}

// PanicParserError panics with a parser error.
func PanicParserError(msg string, l ...int) {
	panic(NewParserError(msg, l...))
}
