package expression

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

func (t *ExpressionHandler) processSymbolExpression(ex ast.SymbolExpression) tf.Var {
	if v, ok := t.st.Vars.Search(ex.Value); ok {
		return v
	}
	errorutils.Abort(errorutils.UnknownVariable, ex.Value)
	return nil
}
