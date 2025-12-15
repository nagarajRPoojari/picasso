package expression

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

func (t *ExpressionHandler) ProcessIndexingExpression(bh *bc.BlockHolder, ex ast.ComputedExpression) tf.Var {
	base := t.ProcessExpression(bh, ex.Member)

	indices := make([]value.Value, 0)
	for _, i := range ex.Indices {
		v := t.ProcessExpression(bh, i)

		casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.INT64), v.Load(bh))
		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)

		indices = append(indices, c.Load(bh))
	}
	v := base.(*tf.Array).LoadByIndex(bh, indices)
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(v.Type())), v)
}
