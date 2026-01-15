package expression

import (
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
)

// processSymbolExpression returns active variable which is in scope.
// Does a search from current scope to upper blocks.
func (t *ExpressionHandler) processSymbolExpression(ex ast.SymbolExpression) tf.Var {
	if v, ok := t.st.Vars.Search(ex.Value); ok {
		return v
	}
	errorutils.Abort(errorutils.UnknownVariable, ex.Value)
	return nil
}
