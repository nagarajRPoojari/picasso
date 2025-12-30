package ast

import (
	"strings"
)

// BlockStatement represents a sequence of statements enclosed in braces.
// In most languages, this construct introduces a new lexical scope.
type BlockStatement struct {
	SourceLoc
	Body []Statement
}

func (BlockStatement) stmt() {}
func (t BlockStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// VariableDeclarationStatement represents the definition of a new name
// in the current scope. It tracks metadata like mutability (Constant),
// visibility (IsStatic), and concurrency hints (IsAtomic).
type VariableDeclarationStatement struct {
	SourceLoc
	Identifier    string
	Constant      bool
	AssignedValue Expression
	ExplicitType  Type
	IsStatic      bool
	IsAtomic      bool
}

func (VariableDeclarationStatement) stmt() {}
func (t VariableDeclarationStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// ExpressionStatement allows an expression to be used where a statement
// is required. The resulting value of the expression is typically discarded.
type ExpressionStatement struct {
	SourceLoc
	Expression Expression
}

func (ExpressionStatement) stmt() {}
func (t ExpressionStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// Parameter represents an input variable definition within a function
// signature, binding a name to a specific type.
type Parameter struct {
	Name string
	Type Type
}

// FunctionDefinitionStatement represents a full function implementation,
// including its signature and executable body. The Hash field can be used
// for memoization or unique identification during IR generation.
type FunctionDefinitionStatement struct {
	SourceLoc
	Parameters []Parameter
	Name       string
	Body       []Statement
	Hash       uint32
	ReturnType Type
	IsStatic   bool
}

func (FunctionDefinitionStatement) stmt() {}
func (t FunctionDefinitionStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// FunctionDeclarationStatement represents a function signature without
// an implementation, often used for interface methods or external linking.
type FunctionDeclarationStatement struct {
	SourceLoc
	Parameters []Parameter
	Name       string
	ReturnType Type
	IsStatic   bool
}

func (FunctionDeclarationStatement) stmt() {}
func (t FunctionDeclarationStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// ReturnStatement terminates the current function execution. It wraps an
// ExpressionStatement to represent the returned value.
type ReturnStatement struct {
	SourceLoc
	Value  ExpressionStatement
	IsVoid bool
}

func (ReturnStatement) stmt() {}
func (t ReturnStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// IfStatement represents conditional branching. The Consequent is executed
// if the Condition is true; otherwise, the optional Alternate is executed.
type IfStatement struct {
	SourceLoc
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

func (IfStatement) stmt() {}
func (t IfStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

// ImportStatement represents a dependency on an external module or file.
// e.g, using "os/bin" as bin;
type ImportStatement struct {
	SourceLoc
	Name  string
	Alias string
}

func (ImportStatement) stmt() {}
func (t ImportStatement) GetSrc() SourceLoc {
	return t.SourceLoc
}

func (t ImportStatement) IsBasePkg() bool {
	return strings.HasPrefix(t.Name, "builtin")
}

func (t ImportStatement) EndName() string {
	pathSplit := strings.Split(t.Name, ".")
	return pathSplit[len(pathSplit)-1]
}

// ForeachStatement represents a collection-based loop. If Index is true,
// the iteration provides the current offset or key.
type ForeachStatement struct {
	SourceLoc
	Value    string
	Index    bool
	Iterable Expression
	Body     []Statement
}

func (n ForeachStatement) stmt() {}
func (n ForeachStatement) GetSrc() SourceLoc {
	return n.SourceLoc
}

// WhileStatement represents a basic loop that continues as long as the
// Condition evaluates to true.
type WhileStatement struct {
	SourceLoc
	Condition Expression
	Body      []Statement
}

func (n WhileStatement) stmt() {}
func (n WhileStatement) GetSrc() SourceLoc {
	return n.SourceLoc
}

// ClassDeclarationStatement represents a blueprint for object instantiation,
// defining encapsulated state (fields) and behavior (methods).
type ClassDeclarationStatement struct {
	SourceLoc
	Name       string
	Body       []Statement
	Implements string
}

func (n ClassDeclarationStatement) stmt() {}
func (n ClassDeclarationStatement) GetSrc() SourceLoc {
	return n.SourceLoc
}

// InterfaceDeclarationStatement defines a contract of method signatures
// that concrete classes must implement.
type InterfaceDeclarationStatement struct {
	SourceLoc
	Name string
	Body []Statement
}

func (n InterfaceDeclarationStatement) stmt() {}
func (n InterfaceDeclarationStatement) GetSrc() SourceLoc {
	return n.SourceLoc
}

// BreakStatement represents an immediate exit from the innermost
// looping construct (For, While, or Foreach).
type BreakStatement struct {
	SourceLoc
}

func (n BreakStatement) stmt() {}
func (n BreakStatement) GetSrc() SourceLoc {
	return n.SourceLoc
}
