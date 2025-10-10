package expression

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *ExpressionHandler) ProcessIndexingExpression(bh tf.BlockHolder, ex ast.ComputedExpression) (tf.Var, tf.BlockHolder) {
	base, safe := t.ProcessExpression(bh, ex.Member)
	bh = safe

	indices := make([]value.Value, 0)
	for _, i := range ex.Indices {
		v, safe := t.ProcessExpression(bh, i)
		bh = safe

		casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, string(tf.INT64), v.Load(bh.N))
		bh.N = safeN
		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)

		indices = append(indices, c.Load(bh.N))
	}
	v := base.(*tf.Array).LoadByIndex(bh.N, indices)
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(v.Type())), v), bh
}
