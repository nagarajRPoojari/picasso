package statement

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

func (t *StatementHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	return expression.ExpressionHandlerInst.ProcessNewExpression(bh, ex)
}
