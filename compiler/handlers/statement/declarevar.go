package statement

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// DeclareVariable handles variable declarations, optionally initializing
// them with an assigned value.
//
// Parameters:
//
//	block - the current IR block
//	st    - the AST VariableDeclarationStatement node
//
// Returns:
//
//	*ir.Block - the updated IR block after declaration
func (t *StatementHandler) DeclareVariable(bh tf.BlockHolder, st *ast.VariableDeclarationStatement) tf.BlockHolder {
	if t.st.Vars.Exists(st.Identifier) {
		errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
	}

	var v tf.Var
	if st.AssignedValue == nil {
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), nil)
	} else {
		_v, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)
		v = _v
		bh = safe
		casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, st.ExplicitType.Get(), v.Load(bh.N))
		bh.N = safeN
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), casted)
	}
	t.st.Vars.AddNewVar(st.Identifier, v)
	return bh
}
