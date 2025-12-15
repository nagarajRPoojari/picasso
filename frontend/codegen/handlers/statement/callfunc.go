package statement

import (
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

func (t *StatementHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {
	return expression.ExpressionHandlerInst.CallFunc(bh, ex)
}
