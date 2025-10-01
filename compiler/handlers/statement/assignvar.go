package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// AssignVariable handles assignment statements, updating either a variable or
// a class member with the result of an evaluated expression.
//
// Parameters:
//
//	block - the current IR block
//	st    - the AST AssignmentExpression node
//
// Returns:
//
//	*ir.Block - the updated IR block after performing the assignment
func (t *StatementHandler) AssignVariable(block *ir.Block, st *ast.AssignmentExpression) *ir.Block {

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		assignee := m.Value
		v, ok := t.st.Vars.Search(assignee)
		if !ok {
			errorutils.Abort(errorutils.UnknownVariable, st)
		}
		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		block = safe

		typeName := v.NativeTypeString()

		casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, typeName, rhs.Load(block))
		block = safe

		c := t.st.TypeHandler.BuildVar(block, tf.Type(typeName), casted)
		v.Update(block, c.Load(block))

	case ast.MemberExpression:

		baseVar, safe := expression.ExpressionHandlerInst.ProcessExpression(block, m.Member)
		block = safe

		if baseVar == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "nil base for member expression")
		}

		// Base must be a class instance
		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "member access base is not a class instance")
		}

		classMeta := t.st.Classes[cls.Name]
		structType := classMeta.StructType()
		meta := t.st.Classes[cls.Name]
		fqName := t.st.IdentifierBuilder.Attach(cls.Name, m.Property)
		index := meta.FieldIndexMap[fqName]

		fieldType := structType.Fields[index]

		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		block = safe

		typeName := fieldType.Name()
		if typeName == "" {
			typeName = fieldType.String()
		}
		if typeName[0:1] == "%" {
			typeName = typeName[1 : len(typeName)-1]
		}
		casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, typeName, rhs.Load(block))
		block = safe
		c := t.st.TypeHandler.BuildVar(block, tf.Type(typeName), casted)
		cls.UpdateField(block, index, c.Load(block), fieldType)
	}

	return block
}
