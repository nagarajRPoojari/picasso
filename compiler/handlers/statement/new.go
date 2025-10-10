package statement

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

func (t *StatementHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	return expression.ExpressionHandlerInst.ProcessNewExpression(bh, ex)
}
