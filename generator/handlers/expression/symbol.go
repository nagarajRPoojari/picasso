package expression

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
)

func (t *ExpressionHandler) processSymbolExpression(ex ast.SymbolExpression) tf.Var {
	if v, ok := t.st.Vars.Search(ex.Value); ok {
		return v
	}
	errorutils.Abort(errorutils.UnknownVariable, ex.Value)
	return nil
}
