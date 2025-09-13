package ast

import "github.com/nagarajRPoojari/x-lang/lexer"

type NumberExpression struct {
	Value float64
}

func (NumberExpression) expr() {
}

type StringExpression struct {
	Value string
}

func (StringExpression) expr() {
}

type SymbolExpression struct {
	Value string
}

func (SymbolExpression) expr() {
}

type BinaryExpression struct {
	Left     Expression
	Operator lexer.Token
	Right    Expression
}

func (BinaryExpression) expr() {}

type AssignmentExpression struct {
	Assignee      Expression
	AssignedValue Expression
}

func (AssignmentExpression) expr() {}

type PrefixExpression struct {
	Operator lexer.Token
	Operand  Expression
}

func (PrefixExpression) expr() {}

type MemberExpression struct {
	Member   Expression
	Property string
}

func (n MemberExpression) expr() {}

type CallExpression struct {
	Method    Expression
	Arguments []Expression
}

func (CallExpression) expr() {}

type ComputedExpression struct {
	Member   Expression
	Property Expression
}

func (ComputedExpression) expr() {}

type RangeExpression struct {
	Lower Expression
	Upper Expression
}

func (RangeExpression) expr() {}

type FunctionExpression struct {
	Parameters []Parameter
	Body       []Statement
	ReturnType Type
}

func (FunctionExpression) expr() {}

type ListExpression struct {
	Constants []Expression
}

func (ListExpression) expr() {}

type NewExpression struct {
	Instantiation CallExpression
}

func (NewExpression) expr() {}
