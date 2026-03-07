package statement

import (
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

// DeclareVariable handles variable declarations, optionally initializing
// them with an assigned value.
// Key Logic:
//   - For unassigned vars initialize then with corresponding zero value, except atomic types.
//   - Do implicit typecasting for assigned vars.
//   - Supports multiple variable declarations: say a: int, b: int = 100, 200;
func (t *StatementHandler) DeclareVariable(bh *bc.BlockHolder, st *ast.VariableDeclarationStatement) {
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)

	// Normalize to multiple declaration format
	identifiers := st.Identifiers
	explicitTypes := st.ExplicitTypes
	var assignedValues []ast.Expression

	if len(st.Identifiers) == 0 {
		identifiers = []string{st.Identifier}
		explicitTypes = []ast.Type{st.ExplicitType}
		if st.AssignedValue != nil {
			assignedValues = []ast.Expression{st.AssignedValue}
		}
	} else {
		assignedValues = st.AssignedValues
	}

	// Check for redeclarations
	for _, id := range identifiers {
		if t.st.Vars.Exists(id) {
			errorutils.Abort(errorutils.VariableRedeclaration, id)
		}
	}

	// if rhs is provided then must provide for all lhs types: sucessfull tuple unpacking
	var declaredVars []tf.Var
	if len(assignedValues) > 0 {
		rhsVars := t.processRHSValues(bh, expHandler, assignedValues, len(identifiers))
		for i := range identifiers {
			declaredVars = append(declaredVars, t.buildTypedVar(bh, explicitTypes[i], rhsVars[i]))
		}
	} else {
		// zero inits
		for i := range identifiers {
			declaredVars = append(declaredVars, t.buildZeroVar(bh, explicitTypes[i]))
		}
	}

	for i, id := range identifiers {
		t.st.Vars.AddNewVar(id, declaredVars[i])
	}
}

// buildTypedVar creates a typed variable from an RHS value with implicit type casting
func (t *StatementHandler) buildTypedVar(bh *bc.BlockHolder, explicitType ast.Type, rhsVar tf.Var) tf.Var {
	tp := t.st.ResolveAlias(explicitType.Get())
	utp := t.st.ResolveAlias(explicitType.GetUnderlyingType())

	if tp == "array" {
		if arr, ok := rhsVar.(*tf.Array); ok {
			return arr
		}
	}

	casted := t.st.TypeHandler.ImplicitTypeCast(bh, tp, rhsVar.Load(bh))
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(tp, utp), casted)
}

// buildZeroVar creates a zero-initialized variable of the given type
// except atomic vars which are initialized by default
func (t *StatementHandler) buildZeroVar(bh *bc.BlockHolder, explicitType ast.Type) tf.Var {
	tp := t.st.ResolveAlias(explicitType.Get())
	utp := t.st.ResolveAlias(explicitType.GetUnderlyingType())

	var init value.Value
	if explicitType.IsAtomic() {
		// atomic data types are special class types & are not expected to be initialized with
		// new keyword. e.g, say x: atomic int; should do the instantiaion job though it is just
		// a declaration. therefore instantiate with NewClass.
		meta := t.st.Classes[tp]
		c := tf.NewClass(bh, tp, meta.UDT)
		init = c.Load(bh)
	}
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(tp, utp), init)
}
