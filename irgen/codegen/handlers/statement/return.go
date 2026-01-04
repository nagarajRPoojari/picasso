package statement

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// Return handles a return statement by evaluating the expression,
// performing implicit type casting to the function's return type, and emitting a return in the IR.
func (t *StatementHandler) Return(block *bc.BlockHolder, st *ast.ReturnStatement, rt ast.Type) {
	v := t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(block, st.Value.Expression)

	val := v.Load(block)

	if rt == nil {
		block.N.NewRet(nil)
	} else {
		r := t.st.TypeHandler.ImplicitTypeCast(block, rt.Get(), val)
		block.N.NewRet(r)
	}
}
