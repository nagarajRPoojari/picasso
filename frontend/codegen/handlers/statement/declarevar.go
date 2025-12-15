package statement

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	errorutils "github.com/nagarajRPoojari/niyama/frontend/codegen/error"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
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
func (t *StatementHandler) DeclareVariable(bh *bc.BlockHolder, st *ast.VariableDeclarationStatement) {
	if t.st.Vars.Exists(st.Identifier) {
		errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
	}

	var v tf.Var
	if st.AssignedValue == nil {
		var init value.Value
		if st.ExplicitType.IsAtomic() {
			meta := t.st.Classes[st.ExplicitType.Get()]
			c := tf.NewClass(bh, st.ExplicitType.Get(), meta.UDT)
			init = c.Load(bh)
		}
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), init)
	} else {
		_v := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)
		v = _v
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, st.ExplicitType.Get(), v.Load(bh))
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), casted)

	}
	t.st.Vars.AddNewVar(st.Identifier, v)
}
