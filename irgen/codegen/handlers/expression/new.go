package expression

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

func buildAliasNameFromMemExp(m ast.Expression) (string, string) {
	switch mc := m.(type) {
	case ast.SymbolExpression:
		return mc.Value, mc.Value
	case ast.MemberExpression:
		res, _ := buildAliasNameFromMemExp(mc.Member)
		return fmt.Sprintf("%s.%s", res, mc.Property), mc.Property
	}
	return "", ""
}

// callConstructor executes the class constructor immediately following allocation.
// In the Niyama object model, constructors are stored as function pointers within
// the class struct itself. This method retrieves that pointer, prepares the
// user-provided arguments, and injects the 'this' pointer to initialize the
// instance's internal state.
//
// Technical Logic:
//   - Symbol Resolution: Maps the class name to its constructor symbol (usually
//     mangled as ClassName.ClassName) via the IdentifierBuilder.
//   - Dynamic Dispatch: Loads the constructor function pointer from the instance's
//     allocated memory using its pre-calculated struct index.
//   - Argument Marshalling: Evaluates constructor arguments and performs implicit
//     casting to ensure binary compatibility with the LLVM function signature.
//   - Instance Binding: Follows the method ABI by appending the allocated
//     instance pointer as the final 'hidden' argument to the call.
func (t *ExpressionHandler) callConstructor(bh *bc.BlockHolder, cls *tf.Class, ex ast.CallExpression) {
	// Get the method symbol and metadata

	aliasClsName, methodName := buildAliasNameFromMemExp(ex.Method)
	aliasConstructorName := fmt.Sprintf("%s.%s", aliasClsName, methodName)

	meta := t.st.Classes[aliasClsName]

	idx := meta.FieldIndexMap[aliasConstructorName]

	st := meta.StructType()
	fieldType := st.Fields[idx]

	// Load the function pointer directly from the struct field (single load)
	fnVal := cls.LoadField(bh, idx, fieldType)
	if fnVal == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("function pointer is nil for %s.%s", cls.Name, aliasClsName))
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

	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v := t.ProcessExpression(bh, argExp)
		if v == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("nil arg %d for %s", i, aliasClsName))
		}
		raw := v.Load(bh)
		if raw == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("loaded nil arg %d for %s", i, aliasClsName))
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

// ProcessNewExpression orchestrates the lifecycle of a new class instance.
// It performs heap allocation (via tf.NewClass), initializes the internal
// struct fields with default values or initializers defined in the AST,
// and finally executes the class constructor.
//
// Technical Logic:
//   - Memory Setup: Allocates the underlying LLVM struct and manages a
//     temporary function-level variable scope for the initialization phase.
//   - Field Initialization: Iterates through the MetaClass field map to
//     differentiate between data fields (variables) and method pointers.
//   - Recursive Type Support: Handles atomic vs. complex type initialization
//     and performs implicit type casting for assigned initial values.
//   - Constructor Dispatch: Finalizes the object state by calling the
//     corresponding constructor method with the 'this' pointer
func (t *ExpressionHandler) ProcessNewExpression(bh *bc.BlockHolder, ex ast.NewExpression) tf.Var {
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	aliasClsName, _ := buildAliasNameFromMemExp(ex.Instantiation.Method)

	if _, ok := t.st.Interfaces[aliasClsName]; ok {
		errorutils.Abort(errorutils.InterfaceInstantiationError, aliasClsName)
	}

	classMeta := t.st.Classes[aliasClsName]
	if classMeta == nil {
		errorutils.Abort(errorutils.UnknownClass, aliasClsName)
	}

	clsNameSplit := strings.Split(aliasClsName, ".")
	moduleName := strings.Join(clsNameSplit[:len(clsNameSplit)-1], ".")

	if classMeta.Internal && moduleName != t.st.ModuleName {
		errorutils.Abort(errorutils.ClassNotAccessible, aliasClsName)
	}

	// tf.NewClass allocates memory for class instance in heap internally.
	// & holds heap pointer in a stack slot.
	instance := tf.NewClass(bh, aliasClsName, classMeta.UDT)
	structType := classMeta.StructType()
	meta := t.st.Classes[aliasClsName]

	for name, index := range meta.FieldIndexMap {
		// if field is not found in meta.VarAST indicating func type, update instance
		// to point to the function. future function calls on that instance will directly
		// refer to this pointed function.
		exp, ok := meta.VarAST[name]
		if !ok {
			f := t.st.Classes[aliasClsName].Methods[name]
			fieldType := t.st.Classes[aliasClsName].StructType().Fields[index]
			instance.UpdateField(bh, t.st.TypeHandler, index, f, fieldType)
			continue
		}

		// handling initialized & uninitialized variables.
		fieldType := structType.Fields[index]
		var v tf.Var
		if exp.AssignedValue == nil {
			// @todo: this need to be verified
			var init value.Value

			// atomic data types are special class types & are not expected to be initialized with
			// new keyword. e.g, say x: atomic int; should do the instantiaion job though it is just
			// a declaration. therefore instantiate with NewClass.
			if exp.ExplicitType.IsAtomic() {
				meta := t.st.Classes[exp.ExplicitType.Get()]
				v = tf.NewClass(bh, exp.ExplicitType.Get(), meta.UDT)
			} else {

				// remaining vars without assignedvalues holds its corresponding zero values.
				// @todo: list zero values for all data types somewehere in docs to look at.
				v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get(), exp.ExplicitType.GetUnderlyingType()), init)
			}
		} else {
			v = t.ProcessExpression(bh, exp.AssignedValue)

			// data types other than array, like primitives, object types are typecasted implicitly
			// before assignment.
			if v.NativeTypeString() != constants.ARRAY {
				casted := t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
				v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get()), casted)
			} else {
				// no need to cast array type, but do a base type check.
				t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
			}
		}
		instance.UpdateField(bh, t.st.TypeHandler, index, v.Load(bh), fieldType)
		t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	t.callConstructor(bh, instance, ex.Instantiation)
	return instance
}
