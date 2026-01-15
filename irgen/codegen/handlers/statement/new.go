package statement

import (
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

// ProcessNewExpression simply delegates call to new expression handler ignoring return type (instance)
// in runtime.
func (t *StatementHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	return t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessNewExpression(bh, ex)
}
