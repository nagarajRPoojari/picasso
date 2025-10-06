package expression

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *ExpressionHandler) ProcessIndexingExpression(block *ir.Block, ex ast.ComputedExpression) (tf.Var, *ir.Block) {
	base, safe := t.ProcessExpression(block, ex.Member)
	block = safe

	indices := make([]value.Value, 0)
	for _, i := range ex.Indices {
		v, safe := t.ProcessExpression(block, i)
		block = safe

		casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, string(tf.INT64), v.Load(block))
		block = safe
		c := t.st.TypeHandler.BuildVar(block, tf.NewType(tf.INT64), casted)

		indices = append(indices, c.Load(block))
	}
	v := base.(*tf.Array).LoadByIndex(block, indices)
	return t.st.TypeHandler.BuildVar(block, tf.NewType(utils.GetTypeString(v.Type())), v), block
}
