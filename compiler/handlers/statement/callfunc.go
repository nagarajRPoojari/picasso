package statement

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

func (t *StatementHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {
	return expression.ExpressionHandlerInst.CallFunc(bh, ex)
}
