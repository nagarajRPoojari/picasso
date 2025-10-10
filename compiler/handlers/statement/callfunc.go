package statement

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *StatementHandler) CallFunc(bh tf.BlockHolder, ex ast.CallExpression) (tf.Var, tf.BlockHolder) {
	return expression.ExpressionHandlerInst.CallFunc(bh, ex)
}
