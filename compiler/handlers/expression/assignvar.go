package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *ExpressionHandler) AssignVariable(block *ir.Block, st *ast.AssignmentExpression) {

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		assignee := m.Value
		v, ok := t.st.Vars.Search(assignee)
		if !ok {
			panic(fmt.Sprintf("undefined: %s", st))
		}
		rhs := t.ProcessExpression(block, st.AssignedValue)
		typeName := v.Type().Name()
		if typeName == "" {
			typeName = v.Type().String()
		}
		casted := t.st.TypeHandler.CastToType(block, typeName, rhs.Load(block))
		c := t.st.TypeHandler.BuildVar(block, tf.Type(typeName), casted)
		v.Update(block, c.Load(block))

	case ast.MemberExpression:
		lhs := t.ProcessExpression(block, m)
		rhs := t.ProcessExpression(block, st.AssignedValue)
		lhs.Update(block, rhs.Load(block))
	}

}
