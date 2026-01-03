package statement

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// CallFunc statements just uses function call expression handler ignoring return type of
// function in runtime.
func (t *StatementHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)
	return expHandler.CallFunc(bh, ex)
}
