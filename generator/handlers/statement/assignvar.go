package statement

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/constants"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/expression"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
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
func (t *StatementHandler) AssignVariable(bh *bc.BlockHolder, st *ast.AssignmentExpression) {

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		assignee := m.Value
		v, ok := t.st.Vars.Search(assignee)
		if !ok {
			errorutils.Abort(errorutils.UnknownVariable, st)
		}
		rhs := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)

		if v.NativeTypeString() != constants.ARRAY {
			typeName := v.NativeTypeString()

			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))

			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
			v.Update(bh, rhs.Load(bh))
		} else {
			v.(*tf.Array).UpdateV2(bh, rhs.(*tf.Array))
		}

	case ast.MemberExpression:

		baseVar := expression.ExpressionHandlerInst.ProcessExpression(bh, m.Member)

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

		rhs := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)

		typeName := utils.GetTypeString(fieldType)
		if typeName != constants.ARRAY {
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
		}
		cls.UpdateField(bh, index, rhs.Load(bh), fieldType)

	case ast.ComputedExpression:
		base := expression.ExpressionHandlerInst.ProcessExpression(bh, m.Member)
		indices := make([]value.Value, 0)
		for _, i := range m.Indices {
			v := expression.ExpressionHandlerInst.ProcessExpression(bh, i)
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.INT64), v.Load(bh))
			c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)
			indices = append(indices, c.Load(bh))
		}

		rhs := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)

		needed := base.(*tf.Array).ElemType

		casted := t.st.TypeHandler.ImplicitTypeCast(bh, utils.GetTypeString(needed), rhs.Load(bh))

		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(needed)), casted)
		base.(*tf.Array).StoreByIndex(bh, indices, c.Load(bh))

	}
}
