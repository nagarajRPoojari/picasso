package statement

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// Return handles a return statement by evaluating the expression,
// performing implicit type casting to the function's return type, and emitting a return in the IR.
//
// Parameters:
//
//	block - the current IR block
//	st    - the AST ReturnStatement node
//	rt    - the expected return type of the function
func (t *StatementHandler) Return(block *bc.BlockHolder, st *ast.ReturnStatement, rt types.Type) {
	v := expression.ExpressionHandlerInst.ProcessExpression(block, st.Value.Expression)

	val := v.Load(block)
	tp := utils.GetTypeString(rt)

	r := t.st.TypeHandler.ImplicitTypeCast(block, tp, val)
	block.N.NewRet(r)
}
