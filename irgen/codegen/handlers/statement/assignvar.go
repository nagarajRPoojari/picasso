package statement

import (
	"fmt"

	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// AssignVariable serves as the central dispatcher for updating storage locations
// within the LLVM IR, handling various forms of Left-Hand Side (LHS) expressions.
// It bridges the gap between high-level language assignment syntax and specific
// memory-access patterns (Load/Store) required by the target architecture.
//
// Technical Logic:
//   - Dispatching: Identifies the assignment target type via AST switching,
//     supporting direct symbols, member access (dot notation), and computed array indices.
//   - Type Coercion: Automatically invokes the ImplicitTypeCast handler to ensure
//     the Right-Hand Side (RHS) matches the storage type before memory updates.
//   - Class Member Resolution: For MemberExpressions, it calculates memory offsets
//     by mapping field names to structural indices via class metadata.
//   - Multi-Dimensional Array Logic: For ComputedExpressions, it resolves nested
//     index access and generates the appropriate element-addressing instructions.
//   - State Persistence: Synchronizes the IR state by calling .Update() or
//     .StoreByIndex() on the variable containers to reflect the new values.
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
			// do implicit type casting to lhs type.
			typeName := v.NativeTypeString()
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
			v.Update(bh, rhs.Load(bh))
		} else {
			// special case to avoid any new var allocation.
			// simply updates pointer to point to new address.
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
		fqName := fmt.Sprintf("%s.%s", cls.Name, m.Property)
		index := meta.FieldIndexMap[fqName]

		fieldType := structType.Fields[index]

		rhs := expression.ExpressionHandlerInst.ProcessExpression(bh, st.AssignedValue)

		typeName := utils.GetTypeString(fieldType)

		// similar to above logic, avoid casting & new var creation logic for array types
		if typeName != constants.ARRAY {
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
		}
		cls.UpdateField(bh, t.st.TypeHandler, index, rhs.Load(bh), fieldType)

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
