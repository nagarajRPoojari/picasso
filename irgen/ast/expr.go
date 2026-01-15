package ast

import "github.com/nagarajRPoojari/picasso/irgen/lexer"

// NumberExpression represents a numeric literal in the source code.
// It stores all numeric values as float64 to maintain high precision
// before type-lowering in the semantic phase.
type NumberExpression struct {
	SourceLoc
	Value float64
}

func (NumberExpression) expr() {}
func (t NumberExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// StringExpression represents a literal sequence of characters.
// This node contains the raw or processed string value, excluding the
// surrounding delimiters (quotes).
type StringExpression struct {
	SourceLoc
	Value string
}

func (StringExpression) expr() {}
func (t StringExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// SymbolExpression represents an identifier used as an expression.
// This typically refers to a variable, constant, or function name
// that must be resolved in the symbol table.
type SymbolExpression struct {
	SourceLoc
	Value string
}

func (SymbolExpression) expr() {}
func (t SymbolExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// BinaryExpression represents an operation involving two operands and an operator.
// It is used for arithmetic, bitwise, and logical operations.
// The Operator token determines the operation's precedence and associativity.
type BinaryExpression struct {
	SourceLoc
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (BinaryExpression) expr() {}
func (t BinaryExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// AssignmentExpression represents the binding of a value to a memory location.
// Invariants: The Assignee must be a valid l-value (e.g., SymbolExpression,
// MemberExpression, or ComputedExpression).
type AssignmentExpression struct {
	SourceLoc
	Assignee      Expression
	AssignedValue Expression
}

func (AssignmentExpression) expr() {}
func (t AssignmentExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// PrefixExpression represents a unary operator that precedes its operand.
// Common examples include logical negation (!), unary minus (-), or pointer dereference.
type PrefixExpression struct {
	SourceLoc
	Operator lexer.Token
	Operand  Expression
}

func (PrefixExpression) expr() {}
func (t PrefixExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// MemberExpression represents static property access on an object or structure
// using dot notation (e.g., object.property).
type MemberExpression struct {
	SourceLoc
	Member   Expression
	Property string
}

func (n MemberExpression) expr() {}
func (t MemberExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// CallExpression represents the invocation of a function, method, or closure.
// The Method expression is evaluated first to determine the callable target.
type CallExpression struct {
	SourceLoc
	Method    Expression
	Arguments []Expression
}

func (CallExpression) expr() {}
func (t CallExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// ComputedExpression represents dynamic access into a collection or structure
// using index notation (e.g., array[index]).
type ComputedExpression struct {
	SourceLoc
	Member  Expression
	Indices []Expression
}

func (ComputedExpression) expr() {}
func (t ComputedExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// RangeExpression represents a span of values, typically defined by a start
// and end point (e.g., 1..10). This is often used in loop iterators.
type RangeExpression struct {
	SourceLoc
	Lower Expression
	Upper Expression
}

func (RangeExpression) expr() {}
func (t RangeExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// FunctionExpression represents an anonymous function or lambda definition.
// It captures the signature and the executable block, allowing functions
// to be treated as first-class citizens.
type FunctionExpression struct {
	SourceLoc
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}

func (FunctionExpression) expr() {}
func (t FunctionExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// ListExpression represents a literal array or list initialization.
// Example: [1, 2, 3 + x, 4].
type ListExpression struct {
	SourceLoc
	Constants []Expression
}

func (ListExpression) expr() {}
func (t ListExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// NewExpression represents the instantiation of a type, typically allocating
// memory and calling a constructor. It wraps a CallExpression to capture
// the type name and arguments.
type NewExpression struct {
	SourceLoc
	Instantiation CallExpression
}

func (NewExpression) expr() {}
func (t NewExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}

// NullExpression represents the explicit absence of a value or a null pointer literal.
type NullExpression struct {
	SourceLoc
}

func (NullExpression) expr() {}
func (t NullExpression) GetSrc() SourceLoc {
	return t.SourceLoc
}
