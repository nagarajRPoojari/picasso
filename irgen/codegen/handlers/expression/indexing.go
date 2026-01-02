package expression

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// ProcessIndexingExpression handles the IR generation for accessing elements within
// an array or a collection using one or more indices (e.g., arr[i] or matrix[i, j]).
// It resolves the base object, coerces all index expressions to 64-bit integers,
// and calculates the memory offset required to fetch the specific element.
//
// Technical Logic:
//   - Index Standardization: Automatically casts all index values (regardless of
//     original width) to INT64 (i64) to ensure compatibility with LLVM GEP instructions.
//   - Recursive Resolution: Processes both the base expression and each index
//     expression through the global ExpressionHandler.
//   - Offset Calculation: Delegates the actual pointer arithmetic to the Array
//     type's LoadByIndex method, which handles multi-dimensional stride calculations
func (t *ExpressionHandler) ProcessIndexingExpression(bh *bc.BlockHolder, ex ast.ComputedExpression) tf.Var {
	base := t.ProcessExpression(bh, ex.Member)

	indices := make([]value.Value, 0)
	for _, i := range ex.Indices {
		v := t.ProcessExpression(bh, i)

		// may be i can safely replace int64 with something lower dtype. need to check @todo
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.INT64), v.Load(bh))
		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)

		indices = append(indices, c.Load(bh))
	}
	arr := base.(*tf.Array)
	v := arr.LoadByIndex(bh, indices)
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(arr.ElementTypeString), v)
}
