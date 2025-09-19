package expression

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) DeclareVariable(block *ir.Block, st *ast.VariableDeclarationStatement) {
	var v tf.Var
	if st.AssignedValue == nil {
		v = t.st.TypeHandler.BuildVar(block, tf.Type(st.ExplicitType.Get()), nil)
		t.st.Vars.AddNewVar(st.Identifier, v)
	} else {
		v = t.ProcessExpression(block, st.AssignedValue)
		if t.st.Vars.Exists(st.Identifier) {
			errorsx.PanicCompilationError("variable already exists")
		}
		casted := t.st.TypeHandler.CastToType(block, st.ExplicitType.Get(), v.Load(block))
		vv := t.st.TypeHandler.BuildVar(block, tf.Type(st.ExplicitType.Get()), casted)
		t.st.Vars.AddNewVar(st.Identifier, vv)
	}
}
