package statement

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *StatementHandler) AssignVariable(block *ir.Block, st *ast.AssignmentExpression) {

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		assignee := m.Value
		v, ok := t.st.Vars.Search(assignee)
		if !ok {
			panic(fmt.Sprintf("undefined: %s", st))
		}
		rhs := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		typeName := v.NativeTypeString()
		casted := t.st.TypeHandler.CastToType(block, typeName, rhs.Load(block))
		c := t.st.TypeHandler.BuildVar(block, tf.Type(typeName), casted)
		v.Update(block, c.Load(block))

	case ast.MemberExpression:

		baseVar := expression.ExpressionHandlerInst.ProcessExpression(block, m.Member)
		if baseVar == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("nil base in member expression: %v", m))
		}

		// Base must be a class instance
		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("member access base is not a class instance, got %T, while, %v", baseVar, m))
		}

		classMeta := t.st.Classes[cls.Name]
		structType := classMeta.StructType()
		meta := t.st.Classes[cls.Name]
		index := meta.VarIndexMap[m.Property]
		fieldType := structType.Fields[index]

		rhs := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		typeName := fieldType.Name()
		if typeName == "" {
			typeName = fieldType.String()
		}
		if typeName[0:1] == "%" {
			typeName = typeName[1 : len(typeName)-1]
		}
		casted := t.st.TypeHandler.CastToType(block, typeName, rhs.Load(block))
		c := t.st.TypeHandler.BuildVar(block, tf.Type(typeName), casted)
		cls.UpdateField(block, index, c.Load(block), fieldType)
	}

}
