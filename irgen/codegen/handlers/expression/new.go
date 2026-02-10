package expression

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
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
// In the Picasso object model, constructors are stored as function pointers within
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
func (t *ExpressionHandler) callConstructor(bh *bc.BlockHolder, cls *tf.Class, ex ast.CallExpression) value.Value {
	// Get the method symbol and metadata

	aliasClsName, methodName := buildAliasNameFromMemExp(ex.Method)
	aliasConstructorName := fmt.Sprintf("%s.%s", aliasClsName, methodName)

	fqClsName := t.st.ResolveAlias(aliasClsName)
	fqConstructorName := t.st.ResolveAlias(aliasConstructorName)

	meta := t.st.Classes[fqClsName]
	fnVal := t.st.Classes[fqClsName].Methods[fqConstructorName]
	args := t.buildConstructorArgs(bh, fqConstructorName, meta, cls, ex)

	// Call the function pointer
	return bh.N.NewCall(fnVal, args...)
}

// utility function to build constructor args
func (t *ExpressionHandler) buildConstructorArgs(bh *bc.BlockHolder, fqConstructorName string, meta *tf.MetaClass, cls *tf.Class, ex ast.CallExpression) []value.Value {
	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v := t.ProcessExpression(bh, argExp)
		if v == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("nil arg %d for %s", i, fqConstructorName))
		}
		raw := v.Load(bh)
		if raw == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("loaded nil arg %d for %s", i, fqConstructorName))
		}

		// Implicit type cast if needed
		expected := meta.MethodArgs[fqConstructorName][i]
		raw = t.st.TypeHandler.ImplicitTypeCast(bh, t.st.ResolveAlias(expected.Get()), raw)
		if raw == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("ImplicitTypeCast returned nil for arg %d -> %s", i, expected.Get()))
		}
		args = append(args, raw)
	}
	// Append `this` pointer as the last argument
	thisPtr := cls.Load(bh)
	if thisPtr == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalInstantiationError, fmt.Sprintf("this pointer is nil for %s", cls.Name))
	}
	return append(args, thisPtr)
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
	fqClsName := t.st.ResolveAlias(aliasClsName)

	if _, ok := t.st.Interfaces[fqClsName]; ok {
		errorutils.Abort(errorutils.InterfaceInstantiationError, fqClsName)
	}

	classMeta := t.st.Classes[fqClsName]
	if classMeta == nil {
		errorutils.Abort(errorutils.UnknownClass, fqClsName)
	}

	clsNameSplit := strings.Split(fqClsName, ".")
	moduleName := strings.Join(clsNameSplit[:len(clsNameSplit)-1], ".")

	if classMeta.Internal && moduleName != t.st.ModuleName {
		errorutils.Abort(errorutils.ClassNotAccessible, fqClsName)
	}

	// tf.NewClass allocates memory for class instance in heap internally.
	// & holds heap pointer in a stack slot.
	instance := tf.NewClass(bh, fqClsName, classMeta.UDT)
	p := t.callConstructor(bh, instance, ex.Instantiation)
	cls := &tf.Class{
		Name: fqClsName,
		UDT:  classMeta.UDT.(*types.PointerType),
	}

	cls.Update(bh, p)
	return cls
}
