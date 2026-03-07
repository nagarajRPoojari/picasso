package statement

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

// Return handles a return statement by evaluating the expression,
// performing implicit type casting to the function's return type, and emitting a return in the IR.
// For multiple return values, it creates a tuple struct and returns it.
func (t *StatementHandler) Return(block *bc.BlockHolder, st *ast.ReturnStatement, rt ast.Type) {
	expHandler := t.m.GetExpressionHandler().(*expression.ExpressionHandler)

	// Handle multiple return values by constructing
	// wrapper tuple struct & return it
	if len(st.Values) > 0 {
		tupleType, ok := rt.(*ast.TupleType)
		if !ok {
			errorutils.Abort(errorutils.TuplePackFailed, "Multiple return values require tuple return type")
		}

		// Evaluate all return expressions
		returnVars := make([]tf.Var, len(st.Values))
		typeNames := make([]string, len(st.Values))
		tupleFieldTypes := make([]types.Type, len(st.Values))

		for i, expr := range st.Values {
			v := expHandler.ProcessExpression(block, expr)
			expectedType := t.st.ResolveAlias(tupleType.Types[i].Get())
			casted := t.st.TypeHandler.ImplicitTypeCast(block, expectedType, v.Load(block))
			returnVars[i] = t.st.TypeHandler.BuildVar(block, tf.NewType(expectedType), casted)
			typeNames[i] = expectedType
			tupleFieldTypes[i] = t.st.TypeHandler.GetLLVMType(expectedType)
		}

		// Get the registered tuple type from GlobalTypeList using type-based naming
		tupleName := utils.GenerateTupleName(tupleFieldTypes, typeNames)
		registeredType, ok := t.st.GlobalTypeList[tupleName]
		if !ok {
			errorutils.Abort(errorutils.InternalError, fmt.Sprintf("Tuple type %s not found in GlobalTypeList", tupleName))
		}

		structType, ok := registeredType.(*types.StructType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, fmt.Sprintf("Expected struct type for %s, got %T", tupleName, registeredType))
		}

		// Create tuple and return it
		tuple := tf.NewTuple(block, structType, returnVars, typeNames)
		block.N.NewRet(tuple.Load(block))
		return
	}

	// Single return value (backward compatibility)
	v := expHandler.ProcessExpression(block, st.Value.Expression)
	val := v.Load(block)

	if rt == nil {
		block.N.NewRet(nil)
	} else {
		r := t.st.TypeHandler.ImplicitTypeCast(block, t.st.ResolveAlias(rt.Get()), val)
		block.N.NewRet(r)
	}
}
