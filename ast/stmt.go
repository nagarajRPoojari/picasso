package ast

type BlockStatement struct {
	Body []Statement
}

func (BlockStatement) stmt() {}

type VariableDeclarationStatement struct {
	Identifier    string
	Constant      bool
	AssignedValue Expression
	ExplicitType  Type
	IsStatic      bool
}

func (VariableDeclarationStatement) stmt() {}

type ExpressionStatement struct {
	Expression Expression
}

func (ExpressionStatement) stmt() {}

type Parameter struct {
	Name string
	Type Type
}

type FunctionDeclarationStatement struct {
	Parameters []Parameter
	Name       string
	Body       []Statement
	ReturnType Type
	IsStatic   bool
}

func (FunctionDeclarationStatement) stmt() {}

type ReturnStatement struct {
	Value ExpressionStatement
}

func (ReturnStatement) stmt() {}

type IfStatement struct {
	Condition  Expression
	Consequent Statement
	Alternate  Statement
}

func (IfStatement) stmt() {}

type ImportStatement struct {
	Name string
	From string
}

func (ImportStatement) stmt() {}

type ForeachStatement struct {
	Value    string
	Index    bool
	Iterable Expression
	Body     []Statement
}

func (n ForeachStatement) stmt() {}

type ClassDeclarationStatement struct {
	Name string
	Body []Statement
}

func (n ClassDeclarationStatement) stmt() {}
