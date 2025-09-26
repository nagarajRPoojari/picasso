package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
)

func (t *StatementHandler) Return(block *ir.Block, st *ast.ReturnStatement, rt types.Type) {
	v, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.Value.Expression)
	block = safe

	val := v.Load(block)
	tp := utils.GetTypeString(rt)

	r, safe := t.st.TypeHandler.ImplicitTypeCast(block, tp, val)
	block = safe

	block.NewRet(r)
}
