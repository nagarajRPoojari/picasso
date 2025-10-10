package statement

import (
	"github.com/llir/llvm/ir"
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
func (t *StatementHandler) DeclareVariable(block *ir.Block, st *ast.VariableDeclarationStatement) *ir.Block {
	if t.st.Vars.Exists(st.Identifier) {
		errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
	}

	var v tf.Var
	if st.AssignedValue == nil {
		v = t.st.TypeHandler.BuildVar(block, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), nil)
	} else {
		_v, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		v = _v
		block = safe
		casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, st.ExplicitType.Get(), v.Load(block))
		block = safe
		v = t.st.TypeHandler.BuildVar(block, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), casted)
	}
	t.st.Vars.AddNewVar(st.Identifier, v)
	return block
}
