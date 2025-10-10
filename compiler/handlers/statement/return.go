package statement

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// Return handles a return statement by evaluating the expression,
// performing implicit type casting to the function's return type, and emitting a return in the IR.
//
// Parameters:
//
//	block - the current IR block
//	st    - the AST ReturnStatement node
//	rt    - the expected return type of the function
func (t *StatementHandler) Return(block tf.BlockHolder, st *ast.ReturnStatement, rt types.Type) {
	v, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.Value.Expression)
	block = safe

	val := v.Load(block.N)
	tp := utils.GetTypeString(rt)

	r, safeN := t.st.TypeHandler.ImplicitTypeCast(block.N, tp, val)
	block.N = safeN

	block.N.NewRet(r)
}
