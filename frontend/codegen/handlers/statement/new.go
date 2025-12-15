package statement

import (
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

func (t *StatementHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	return expression.ExpressionHandlerInst.ProcessNewExpression(bh, ex)
}
