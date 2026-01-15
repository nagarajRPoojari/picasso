// Package ast declares the types used to represent the Abstract Syntax Tree
// of the Picasso source code. The AST serves as the intermediate representation
// between the Parser and the IR Generator (irgen).
package ast

import "github.com/nagarajRPoojari/picasso/irgen/utils"

type SourceLoc struct {
	FilePath string
	Line     int
	Col      int
}

// Statement represents a node in the AST that performs an action but
// does not evaluate to a value (e.g., assignments, loops, declarations).
type Statement interface {
	// stmt is a dummy method used to ensure type safety, preventing
	// Expression types from being used where Statements are expected.
	stmt()
	GetSrc() SourceLoc
}

// Expression represents a node that evaluates to a specific value at runtime.
// Expressions can be nested to form complex computational logic.
type Expression interface {
	// expr is a dummy method used to differentiate Expression nodes
	// from Statement nodes at compile-time.
	expr()
	GetSrc() SourceLoc
}

// Type defines the interface for the language's type system.
type Type interface {
	// IsAtomic returns true if the type is a atomic type (e.g., atomic int, atomic bool)
	// and cannot be further decomposed.
	IsAtomic() bool
	// SetAtomic marks the type as a atomic type.
	SetAtomic()
	// Get returns the string representation of the type name.
	Get() string
	// GetUnderlyingType returns the base type, resolving composite types like array
	// to its base type.
	GetUnderlyingType() string
}

// ExpectExpr is a generic helper that asserts an Expression is of a specific
// concrete type T. It panics via the utils package if the assertion fails.
//
// Usage:
//
//	lit := ast.ExpectExpr[*ast.Literal](node)
func ExpectExpr[T Expression](expr Expression) T {
	return utils.ExpectType[T](expr)
}

// ExpectStmt is a generic helper that asserts a Statement is of a specific
// concrete type T. It is primarily used during IR generation to narrow
// interface types to concrete AST nodes.
func ExpectStmt[T Statement](expr Statement) T {
	return utils.ExpectType[T](expr)
}
