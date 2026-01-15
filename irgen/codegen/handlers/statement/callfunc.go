package statement

import (
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

// CallFunc statements just uses function call expression handler ignoring return type of
// function in runtime.
func (t *StatementHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)
	return expHandler.CallFunc(bh, ex)
}
