package expression

import (
	"fmt"

	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) processSymbolExpression(ex ast.SymbolExpression) tf.Var {
	if v, ok := t.st.Vars.Search(ex.Value); ok {
		return v
	}
	errorsx.PanicCompilationError(fmt.Sprintf("undefined var: %s", ex.Value))
	return nil
}
