package ast

// BlockStatement represents a sequence of statements enclosed in braces.
// In most languages, this construct introduces a new lexical scope.
type BlockStatement struct {
	Body []Statement
}

func (BlockStatement) stmt() {}

// VariableDeclarationStatement represents the definition of a new name
// in the current scope. It tracks metadata like mutability (Constant),
// visibility (IsStatic), and concurrency hints (IsAtomic).
type VariableDeclarationStatement struct {
	Identifier    string
	Constant      bool
	AssignedValue Expression
	ExplicitType  Type
	IsStatic      bool
	IsAtomic      bool
}

func (VariableDeclarationStatement) stmt() {}

// ExpressionStatement allows an expression to be used where a statement
// is required. The resulting value of the expression is typically discarded.
type ExpressionStatement struct {
	Expression Expression
}

func (ExpressionStatement) stmt() {}

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
	Parameters []Parameter
	Name       string
	Body       []Statement
	Hash       uint32
	ReturnType Type
	IsStatic   bool
}

func (FunctionDefinitionStatement) stmt() {}

// FunctionDeclarationStatement represents a function signature without
// an implementation, often used for interface methods or external linking.
type FunctionDeclarationStatement struct {
	Parameters []Parameter
	Name       string
	ReturnType Type
	IsStatic   bool
}

func (FunctionDeclarationStatement) stmt() {}

// ReturnStatement terminates the current function execution. It wraps an
// ExpressionStatement to represent the returned value.
type ReturnStatement struct {
	Value  ExpressionStatement
	IsVoid bool
}

func (ReturnStatement) stmt() {}

// IfStatement represents conditional branching. The Consequent is executed
// if the Condition is true; otherwise, the optional Alternate is executed.
type IfStatement struct {
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

func (IfStatement) stmt() {}

// ImportStatement represents a dependency on an external module or file.
// 'From' typically represents the path, while 'Name' is the local alias.
type ImportStatement struct {
	Name string
	From string
}

func (ImportStatement) stmt() {}

// ForeachStatement represents a collection-based loop. If Index is true,
// the iteration provides the current offset or key.
type ForeachStatement struct {
	Value    string
	Index    bool
	Iterable Expression
	Body     []Statement
}

func (n ForeachStatement) stmt() {}

// WhileStatement represents a basic loop that continues as long as the
// Condition evaluates to true.
type WhileStatement struct {
	Condition Expression
	Body      []Statement
}

func (n WhileStatement) stmt() {}

// ClassDeclarationStatement represents a blueprint for object instantiation,
// defining encapsulated state (fields) and behavior (methods).
type ClassDeclarationStatement struct {
	Name       string
	Body       []Statement
	Implements string
}

func (n ClassDeclarationStatement) stmt() {}

// InterfaceDeclarationStatement defines a contract of method signatures
// that concrete classes must implement.
type InterfaceDeclarationStatement struct {
	Name string
	Body []Statement
}

func (n InterfaceDeclarationStatement) stmt() {}

// BreakStatement represents an immediate exit from the innermost
// looping construct (For, While, or Foreach).
type BreakStatement struct{}

func (n BreakStatement) stmt() {}
