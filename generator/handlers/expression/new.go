package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/constants"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

func (t *ExpressionHandler) callConstructor(bh *bc.BlockHolder, cls *tf.Class, ex ast.CallExpression) {
	// Get the method symbol and metadata
	m := ex.Method.(ast.SymbolExpression)
	meth := t.st.IdentifierBuilder.Attach(m.Value, m.Value)
	meta := t.st.Classes[cls.Name]
	idx := meta.FieldIndexMap[meth]

	// Get struct and field type
	st := meta.StructType()
	fieldType := st.Fields[idx]

	// Load the function pointer directly from the struct field (single load)
	fnVal := cls.LoadField(bh, idx, fieldType)
	if fnVal == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("function pointer is nil for %s.%s", cls.Name, m.Value))
	}

	// Ensure field type is pointer-to-function
	var funcType *types.FuncType
	if ptrType, ok := fieldType.(*types.PointerType); ok {
		funcType, ok = ptrType.ElemType.(*types.FuncType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("expected pointer-to-function, got pointer to %T", ptrType.ElemType))
		}
	} else {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("expected pointer-to-function type for field, got %T", fieldType))
	}

	// Build arguments
	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v := t.ProcessExpression(bh, argExp)
		if v == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("nil arg %d for %s", i, m.Value))
		}
		raw := v.Load(bh)
		if raw == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("loaded nil arg %d for %s", i, m.Value))
		}

		// Implicit type cast if needed
		expected := funcType.Params[i]
		target := utils.GetTypeString(expected)
		raw = t.st.TypeHandler.ImplicitTypeCast(bh, target, raw)
		if raw == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("ImplicitTypeCast returned nil for arg %d -> %s", i, target))
		}
		args = append(args, raw)
	}

	// Append `this` pointer as the last argument
	thisPtr := cls.Load(bh)
	if thisPtr == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("this pointer is nil for %s", cls.Name))
	}
	args = append(args, thisPtr)

	// Call the function pointer
	bh.N.NewCall(fnVal, args...)
}

// ProcessNewExpression handles the creation of a new class instance (`new` expression).
//
// Steps:
//  1. Adds a new function scope for variables and ensures cleanup via defer.
//  2. Resolves the class metadata for the given instantiation method.
//  3. Allocates a new class instance and initializes its underlying struct type.
//  4. Iterates over all fields defined in the class:
//     - For fields with assigned AST values, evaluates the expression, performs
//     implicit type casting if necessary, and stores the value.
//     - For fields without assigned values, initializes a default variable
//     using the explicit type specified in the AST.
//     - For method fields, updates the instance with the function pointer.
//  5. Registers each initialized variable in the current scope.
//  6. Invokes the class constructor if specified via `callConstructor`.
//
// Parameters:
//
//	block - the current IR block where the new instance and field assignments are emitted
//	ex    - the AST `NewExpression` node containing instantiation details
//
// Returns:
//
//	tf.Var     - the newly created class instance
//	*ir.Block  - the updated IR block after field initialization and constructor call
func (t *ExpressionHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	meth := ex.Instantiation.Method.(ast.SymbolExpression)
	classMeta := t.st.Classes[meth.Value]
	if classMeta == nil {
		errorutils.Abort(errorutils.UnknownClass, meth.Value)
	}

	instance := tf.NewClass(bh, meth.Value, classMeta.UDT)
	structType := classMeta.StructType()
	meta := t.st.Classes[meth.Value]

	for name, index := range meta.FieldIndexMap {
		exp, ok := meta.VarAST[name]
		if !ok {
			f := t.st.Classes[meth.Value].Methods[name]
			fieldType := t.st.Classes[meth.Value].StructType().Fields[index]
			instance.UpdateField(bh, index, f, fieldType)
			continue
		}
		fieldType := structType.Fields[index]

		var v tf.Var
		if exp.AssignedValue == nil {
			var init value.Value
			if exp.ExplicitType.IsAtomic() {
				meta := t.st.Classes[exp.ExplicitType.Get()]
				c := tf.NewClass(bh, exp.ExplicitType.Get(), meta.UDT)
				init = c.Load(bh)
			}
			v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get(), exp.ExplicitType.GetUnderlyingType()), init)
		} else {
			_v := t.ProcessExpression(bh, exp.AssignedValue)
			v = _v

			if v.NativeTypeString() != constants.ARRAY {
				casted := t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
				v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get()), casted)
			} else {
				// no need to cast, but does type check
				t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
			}
		}
		instance.UpdateField(bh, index, v.Load(bh), fieldType)
		t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	t.callConstructor(bh, instance, ex.Instantiation)
	return instance
}
