package statement

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *StatementHandler) ProcessNewExpression(bh tf.BlockHolder, ex ast.NewExpression) (tf.Var, tf.BlockHolder) {
	return expression.ExpressionHandlerInst.ProcessNewExpression(bh, ex)
}
