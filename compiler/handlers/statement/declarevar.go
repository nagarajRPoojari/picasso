package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *StatementHandler) DeclareVariable(block *ir.Block, st *ast.VariableDeclarationStatement) {
	if t.st.Vars.Exists(st.Identifier) {
		errorsx.PanicCompilationError("variable already exists")
	}

	var v tf.Var
	if st.AssignedValue == nil {
		v = t.st.TypeHandler.BuildVar(block, tf.Type(st.ExplicitType.Get()), nil)
	} else {
		v = expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		casted := t.st.TypeHandler.CastToType(block, st.ExplicitType.Get(), v.Load(block))
		v = t.st.TypeHandler.BuildVar(block, tf.Type(st.ExplicitType.Get()), casted)
	}
	t.st.Vars.AddNewVar(st.Identifier, v)
}
