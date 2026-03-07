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

	// Normalize to multiple assignment format
	assignees := st.Assignees
	var assignedValues []ast.Expression

	if len(st.Assignees) == 0 {
		// Single assignment - convert to slice format
		assignees = []ast.Expression{st.Assignee}
		assignedValues = []ast.Expression{st.AssignedValue}
	} else {
		assignedValues = st.AssignedValues
	}

	// Process RHS values (handles both single and tuple returns)
	rhsVars := t.processRHSValues(bh, expHandler, assignedValues, len(assignees))

	// Assign to each LHS
	for i, assignee := range assignees {
		t.assignToTarget(bh, expHandler, assignee, rhsVars[i])
	}
}

// processRHSValues evaluates RHS expressions and handles tuple unpacking
// Returns a slice of Var values ready for assignment
func (t *StatementHandler) processRHSValues(bh *bc.BlockHolder, expHandler *expression.ExpressionHandler, assignedValues []ast.Expression, expectedCount int) []tf.Var {
	var rhsVars []tf.Var

	// Check if RHS is a single expression that returns a tuple
	if len(assignedValues) == 1 {
		rhsVar := expHandler.ProcessExpression(bh, assignedValues[0])

		if tuple, ok := rhsVar.(*tf.Tuple); ok {
			for i := range expectedCount {
				fieldVal := tuple.GetField(bh, i)
				rhsVars = append(rhsVars, t.st.TypeHandler.BuildVar(bh, tf.NewType(tuple.TypeNames[i]), fieldVal))
			}
		} else {
			if expectedCount == 1 {
				rhsVars = append(rhsVars, rhsVar)
			} else {
				errorutils.Abort(errorutils.TupleUnpackFailed, "Cannot assign single non-tuple value to multiple variables")
			}
		}
	} else {
		// Multiple values on RHS - need to handle tuples that may be unpacked
		for _, expr := range assignedValues {
			_v := expHandler.ProcessExpression(bh, expr)

			if tuple, ok := _v.(*tf.Tuple); ok {
				// Unpack tuple fields
				for i := 0; i < len(tuple.NativeType.Fields); i++ {
					fieldVal := tuple.GetField(bh, i)
					rhsVars = append(rhsVars, t.st.TypeHandler.BuildVar(bh, tf.NewType(tuple.TypeNames[i]), fieldVal))
				}
			} else {
				// Single value
				rhsVars = append(rhsVars, _v)
			}
		}

		if len(rhsVars) != expectedCount {
			errorutils.Abort(errorutils.TupleUnpackFailed, fmt.Sprintf("Assignment count mismatch: %d variables, %d values", expectedCount, len(rhsVars)))
		}
	}

	return rhsVars
}

// assignToTarget assigns a value to a specific target (symbol, member, or computed expression)
func (t *StatementHandler) assignToTarget(bh *bc.BlockHolder, expHandler *expression.ExpressionHandler, target ast.Expression, rhs tf.Var) {
	switch m := target.(type) {
	case ast.SymbolExpression:
		v, ok := t.st.Vars.Search(m.Value)
		if !ok {
			errorutils.Abort(errorutils.UnknownVariable, target)
		}

		if v.NativeTypeString() != constants.ARRAY {
			typeName := v.NativeTypeString()
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
			castedVar := t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
			v.Update(bh, castedVar.Load(bh))
		} else {
			v.(*tf.Array).UpdateV2(bh, rhs.(*tf.Array))
		}

	case ast.MemberExpression:
		baseVar := expHandler.ProcessExpression(bh, m.Member)
		if baseVar == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "nil base for member expression")
		}

		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalMemberExprError, "member access base is not a class instance")
		}

		classMeta := t.st.Classes[cls.Name]
		structType := classMeta.StructType()
		fqName := fmt.Sprintf("%s.%s", cls.Name, m.Property)
		index := classMeta.FieldIndexMap[fqName]
		fieldType := structType.Fields[index]

		typeName := t.st.ResolveAlias(classMeta.VarAST[fqName].ExplicitType.Get())

		if resolveRootMember(m) != constants.THIS {
			if _, ok := classMeta.InternalFields[fqName]; ok {
				errorutils.Abort(errorutils.FieldNotAccessible, cls.Name, m.Property)
			}
		}

		if typeName != constants.ARRAY {
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, typeName, rhs.Load(bh))
			rhs = t.st.TypeHandler.BuildVar(bh, tf.NewType(typeName), casted)
		}
		cls.UpdateField(bh, t.st.TypeHandler, index, rhs.Load(bh), fieldType)

	case ast.ComputedExpression:
		base := expHandler.ProcessExpression(bh, m.Member)
		indices := make([]value.Value, 0)
		for _, idx := range m.Indices {
			v := expHandler.ProcessExpression(bh, idx)
			casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.INT64), v.Load(bh))
			c := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT64), casted)
			indices = append(indices, c.Load(bh))
		}

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
}

// utility function to get root name of memeber expression
func resolveRootMember(ex ast.Expression) string {
	switch st := ex.(type) {
	case ast.SymbolExpression:
		return st.Value
	case ast.MemberExpression:
		return resolveRootMember(st.Member)
	case ast.ComputedExpression:
		return resolveRootMember(st.Member)
	}

	panic("something gone wrong")
}
