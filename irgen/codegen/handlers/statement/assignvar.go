package statement

import (
	"fmt"

	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
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
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)

	switch m := st.Assignee.(type) {
	case ast.SymbolExpression:
		t.processLocalVarAssignment(bh, expHandler, m, st)

	case ast.MemberExpression:
		t.processClassFieldAssignment(bh, expHandler, m, st)

	case ast.ComputedExpression:
		t.processArrayFieldAssignment(bh, expHandler, m, st)
	}
}

func (t *StatementHandler) processLocalVarAssignment(bh *bc.BlockHolder, expHandler *expression.ExpressionHandler, m ast.SymbolExpression, st *ast.AssignmentExpression) {
	assignee := m.Value
	v, ok := t.st.Vars.Search(assignee)
	if !ok {
		errorutils.Abort(errorutils.UnknownVariable, st)
	}
	rhs := expHandler.ProcessExpression(bh, st.AssignedValue)

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

}

func (t *StatementHandler) processClassFieldAssignment(bh *bc.BlockHolder, expHandler *expression.ExpressionHandler, m ast.MemberExpression, st *ast.AssignmentExpression) {
	baseVar := expHandler.ProcessExpression(bh, m.Member)

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

	rhs := expHandler.ProcessExpression(bh, st.AssignedValue)

	typeName := t.st.ResolveAlias(classMeta.VarAST[fqName].ExplicitType.Get())

	// similar to above logic, avoid casting & new var creation logic for array types
	if typeName != constants.ARRAY {
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
		rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
	}
	cls.UpdateField(bh, t.st.TypeHandler, index, rhs.Load(bh), fieldType)
}

func (t *StatementHandler) processArrayFieldAssignment(bh *bc.BlockHolder, expHandler *expression.ExpressionHandler, m ast.ComputedExpression, st *ast.AssignmentExpression) {
	base := expHandler.ProcessExpression(bh, m.Member)
	indices := make([]value.Value, 0)
	for _, i := range m.Indices {
		v := expHandler.ProcessExpression(bh, i)
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.INT64), v.Load(bh))
		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)
		indices = append(indices, c.Load(bh))
	}

	rhs := expHandler.ProcessExpression(bh, st.AssignedValue)
	arr := base.(*tf.Array)

	if len(indices) < arr.Rank {
		rhsArray, ok := rhs.(*tf.Array)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "partial array indexing requires array value on RHS")
		}
		arr.StoreSubarrayByIndex(bh, indices, rhsArray)
	} else {
		needed := arr.ElementTypeString
		casted := t.st.TypeHandler.ImplicitTypeCast(bh, needed, rhs.Load(bh))
		c := t.st.TypeHandler.BuildVar(bh, tf.NewType(needed), casted)
		arr.StoreByIndex(bh, indices, c.Load(bh))
	}
}
