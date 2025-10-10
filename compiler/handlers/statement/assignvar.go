package statement

import (
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
func (t *StatementHandler) AssignVariable(bh tf.BlockHolder, st *ast.AssignmentExpression) tf.BlockHolder {

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		assignee := m.Value
		v, ok := t.st.Vars.Search(assignee)
		if !ok {
			errorutils.Abort(errorutils.UnknownVariable, st)
		}
		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)
		bh = safe

		if v.NativeTypeString() != constants.ARRAY {
			typeName := v.NativeTypeString()

			casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, typeName, rhs.Load(bh.N))
			bh.N = safeN

			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
			v.Update(bh.N, rhs.Load(bh.N))
		} else {
			v.(*tf.Array).UpdateV2(bh.N, rhs.(*tf.Array))
		}

	case ast.MemberExpression:

		baseVar, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, m.Member)
		bh = safe

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

		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)
		bh = safe

		typeName := utils.GetTypeString(fieldType)
		if typeName != constants.ARRAY {
			casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, typeName, rhs.Load(bh.N))
			bh.N = safeN
			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
		}
		cls.UpdateField(bh.N, index, rhs.Load(bh.N), fieldType)

	case ast.ComputedExpression:
		base, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, m.Member)
		bh = safe
		indices := make([]value.Value, 0)
		for _, i := range m.Indices {
			v, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, i)
			bh = safe
			casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, string(tf.INT64), v.Load(bh.N))
			bh.N = safeN
			c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)
			indices = append(indices, c.Load(bh.N))
		}

		rhs, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)
		bh = safe

		needed := base.(*tf.Array).ElemType

		casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, utils.GetTypeString(needed), rhs.Load(bh.N))
		bh.N = safeN

		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(needed)), casted)
		base.(*tf.Array).StoreByIndex(bh.N, indices, c.Load(bh.N))

	}

	return bh
}
