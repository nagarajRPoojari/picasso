package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *StatementHandler) ProcessNewExpression(block *ir.Block, ex ast.NewExpression) (tf.Var, *ir.Block) {
	return expression.ExpressionHandlerInst.ProcessNewExpression(block, ex)
}
