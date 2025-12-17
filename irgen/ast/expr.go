package ast

import "github.com/nagarajRPoojari/niyama/irgen/lexer"

// NumberExpression represents a numeric literal in the source code.
// It stores all numeric values as float64 to maintain high precision
// before type-lowering in the semantic phase.
type NumberExpression struct {
	Value float64
}

func (NumberExpression) expr() {}

// StringExpression represents a literal sequence of characters.
// This node contains the raw or processed string value, excluding the
// surrounding delimiters (quotes).
type StringExpression struct {
	Value string
}

func (StringExpression) expr() {}

// SymbolExpression represents an identifier used as an expression.
// This typically refers to a variable, constant, or function name
// that must be resolved in the symbol table.
type SymbolExpression struct {
	Value string
}

func (SymbolExpression) expr() {}

// BinaryExpression represents an operation involving two operands and an operator.
// It is used for arithmetic, bitwise, and logical operations.
// The Operator token determines the operation's precedence and associativity.
type BinaryExpression struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (BinaryExpression) expr() {}

// AssignmentExpression represents the binding of a value to a memory location.
// Invariants: The Assignee must be a valid l-value (e.g., SymbolExpression,
// MemberExpression, or ComputedExpression).
type AssignmentExpression struct {
	Assignee      Expression
	AssignedValue Expression
}

func (AssignmentExpression) expr() {}

// PrefixExpression represents a unary operator that precedes its operand.
// Common examples include logical negation (!), unary minus (-), or pointer dereference.
type PrefixExpression struct {
	Operator lexer.Token
	Operand  Expression
}

func (PrefixExpression) expr() {}

// MemberExpression represents static property access on an object or structure
// using dot notation (e.g., object.property).
type MemberExpression struct {
	Member   Expression
	Property string
}

func (n MemberExpression) expr() {}

// CallExpression represents the invocation of a function, method, or closure.
// The Method expression is evaluated first to determine the callable target.
type CallExpression struct {
	Method    Expression
	Arguments []Expression
}

func (CallExpression) expr() {}

// ComputedExpression represents dynamic access into a collection or structure
// using index notation (e.g., array[index]).
type ComputedExpression struct {
	Member  Expression
	Indices []Expression
}

func (ComputedExpression) expr() {}

// RangeExpression represents a span of values, typically defined by a start
// and end point (e.g., 1..10). This is often used in loop iterators.
type RangeExpression struct {
	Lower Expression
	Upper Expression
}

func (RangeExpression) expr() {}

// FunctionExpression represents an anonymous function or lambda definition.
// It captures the signature and the executable block, allowing functions
// to be treated as first-class citizens.
type FunctionExpression struct {
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}

func (FunctionExpression) expr() {}

// ListExpression represents a literal array or list initialization.
// Example: [1, 2, 3 + x, 4].
type ListExpression struct {
	Constants []Expression
}

func (ListExpression) expr() {}

// NewExpression represents the instantiation of a type, typically allocating
// memory and calling a constructor. It wraps a CallExpression to capture
// the type name and arguments.
type NewExpression struct {
	Instantiation CallExpression
}

func (NewExpression) expr() {}

// NullExpression represents the explicit absence of a value or a null pointer literal.
type NullExpression struct{}

func (NullExpression) expr() {}
