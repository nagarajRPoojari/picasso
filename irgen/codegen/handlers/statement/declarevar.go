package statement

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// DeclareVariable handles variable declarations, optionally initializing
// them with an assigned value.
// Key Logic:
//   - For unassigned vars initialize then with corresponding zero value, except atomic types.
//   - Do implicit typecasting for assigned vars.
func (t *StatementHandler) DeclareVariable(bh *bc.BlockHolder, st *ast.VariableDeclarationStatement) {
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)
	if t.st.Vars.Exists(st.Identifier) {
		errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
	}

	var v tf.Var
	if st.AssignedValue == nil {
		var init value.Value
		if st.ExplicitType.IsAtomic() {
			// atomic data types are special class types & are not expected to be initialized with
			// new keyword. e.g, say x: atomic int; should do the instantiaion job though it is just
			// a declaration. therefore instantiate with NewClass.
			meta := t.st.Classes[st.ExplicitType.Get()]
			c := tf.NewClass(bh, st.ExplicitType.Get(), meta.UDT)
			init = c.Load(bh)
		}
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), init)
	} else {
		_v := expHandler.ProcessExpression(bh, st.AssignedValue)
		v = _v
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, st.ExplicitType.Get(), v.Load(bh))
		v = t.st.TypeHandler.BuildVar(bh, tf.NewType(st.ExplicitType.Get(), st.ExplicitType.GetUnderlyingType()), casted)

	}
	t.st.Vars.AddNewVar(st.Identifier, v)
}
