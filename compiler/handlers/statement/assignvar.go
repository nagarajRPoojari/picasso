package statement

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
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

		if v.NativeTypeString() != constants.ARRAY {
			typeName := v.NativeTypeString()

			casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, typeName, rhs.Load(block))
			block = safe

			rhs = t.st.TypeHandler.BuildVar(block, tf.NewType(typeName), casted)
			v.Update(block, rhs.Load(block))
		} else {
			v.(*tf.Array).UpdateV2(block, rhs.(*tf.Array))
		}

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

		typeName := utils.GetTypeString(fieldType)
		if typeName != constants.ARRAY {
			casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, typeName, rhs.Load(block))
			block = safe
			rhs = t.st.TypeHandler.BuildVar(block, tf.NewType(typeName), casted)
		}
		cls.UpdateField(block, index, rhs.Load(block), fieldType)

	case ast.ComputedExpression:
		base, safe := expression.ExpressionHandlerInst.ProcessExpression(block, m.Member)
		block = safe
		indices := make([]value.Value, 0)
		for _, i := range m.Indices {
			v, safe := expression.ExpressionHandlerInst.ProcessExpression(block, i)
			block = safe
			casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, string(tf.INT64), v.Load(block))
			block = safe
			c := t.st.TypeHandler.BuildVar(block, tf.NewType(tf.INT64), casted)
			indices = append(indices, c.Load(block))
		}

		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(block, st.AssignedValue)
		block = safe

		needed := base.(*tf.Array).ElemType

		casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, utils.GetTypeString(needed), rhs.Load(block))
		block = safe

		c := t.st.TypeHandler.BuildVar(block, tf.NewType(utils.GetTypeString(needed)), casted)
		base.(*tf.Array).StoreByIndex(block, indices, c.Load(block))
	}

	return block
}
