package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *StatementHandler) CallFunc(block *ir.Block, ex ast.CallExpression) (tf.Var, *ir.Block) {
	return expression.ExpressionHandlerInst.CallFunc(block, ex)
}
