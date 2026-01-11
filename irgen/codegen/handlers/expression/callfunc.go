package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs/libutils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// CallFunc orchestrates function and method invocation in the LLVM IR.
// It resolves function symbols across three primary domains:
// 1. External Library/Imported Methods (Namespace-prefixed)
// 2. Global/Built-in Symbols (Direct symbols or runtime intrinsics)
// 3. Class Instance Methods (Dynamic dispatch via function pointers in structs)
//
// Technical Logic:
//   - Method Resolution: For class members, it fetches the function pointer from
//     the struct layout using GEP (GetElementPtr) via LoadField.
//   - Calling Convention: Implements the "hidden this" pattern by appending the
//     instance pointer as the final argument to method calls.
//   - Type Safety: Performs implicit type casting of arguments to match the
//     formal parameters defined in the function signature.
func (t *ExpressionHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {

	// check imported base modules for method resolution
	if ret, ok := t.callLibMethod(bh, ex); ok {
		return ret
	}

	if ret, ok := t.callFFIMethod(bh, ex); ok {
		return ret
	}

	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		return t.callNativeMethods(bh, ex, m)

	case ast.MemberExpression:
		return t.callClassMethod(bh, ex, m)
	}
	return nil
}

func (t *ExpressionHandler) callLibMethod(bh *bc.BlockHolder, ex ast.CallExpression) (tf.Var, bool) {
	m, ok := ex.Method.(ast.MemberExpression)
	if !ok {
		return nil, false
	}

	x, ok := m.Member.(ast.SymbolExpression)
	if !ok {
		return nil, false
	}

	// base module func calls are strictly expected to be in module_name.func() format. e.g, syncio.printf
	fName := fmt.Sprintf("%s.%s", t.st.Imports[x.Value].Name, m.Property)
	f, ok := t.st.LibMethods[fName]
	if !ok {
		return nil, false
	}

	args := make([]tf.Var, 0)
	for _, v := range ex.Arguments {
		res := t.ProcessExpression(bh, v)
		args = append(args, res)
	}
	ret := f(t.st.TypeHandler, t.st.Module, bh, args)
	return ret, true
}

func (t *ExpressionHandler) callFFIMethod(bh *bc.BlockHolder, ex ast.CallExpression) (tf.Var, bool) {
	m, ok := ex.Method.(ast.MemberExpression)
	if !ok {
		return nil, false
	}

	x, ok := m.Member.(ast.SymbolExpression)
	if !ok {
		return nil, false
	}

	moduleName := t.st.Imports[x.Value].Name
	fnName := m.Property

	ffiModule, ok := t.st.FFIModules[moduleName]
	if !ok {
		return nil, false
	}

	fn, ok := ffiModule.Methods[fnName]
	if !ok {
		return nil, false
	}

	args := make([]tf.Var, 0)
	for _, v := range ex.Arguments {
		res := t.ProcessExpression(bh, v)
		args = append(args, res)
	}

	ret := libutils.CallCFunc(t.st.TypeHandler, fn, bh, args)
	return ret, true
}

func (t *ExpressionHandler) callNativeMethods(bh *bc.BlockHolder, ex ast.CallExpression, methodName ast.SymbolExpression) tf.Var {
	// @fix: refactor this piece of code
	// thread() special functions comes from base modules but exception for module_name.func() format
	// in future i might add many such special funcs, so need to be kept somewhere else.
	if methodName.Value != c.FUNC_THREAD {
		errorutils.Abort(errorutils.UnknownMethod, methodName.Value)
		return nil
	}

	meth, ok := t.st.CI.Funcs[methodName.Value]
	if !ok {
		errorutils.Abort(errorutils.UnknownMethod, methodName.Value)
		return nil
	}

	// Total elements to pass to thread():
	// [func_ptr, nargs, arg1, arg2, ..., this]
	m := ex.Arguments[0].(ast.MemberExpression)
	cls := t.ProcessExpression(bh, ex.Arguments[0].(ast.MemberExpression).Member).(*tf.Class)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "member access base is not a Class type")
	}
	if cls == nil || cls.Ptr == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "class or class.Ptr is nil for class")
	}

	classMeta := t.st.Classes[cls.Name]
	if classMeta == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "unknown class metadata: "+cls.Name)
	}

	methodFqName := fmt.Sprintf("%s.%s", cls.Name, m.Property)
	idx, ok := classMeta.FieldIndexMap[methodFqName]

	if !ok {
		errorutils.Abort(errorutils.UnknownMethod, m.Property)
	}

	st := classMeta.StructType()
	fieldType := st.Fields[idx]

	// Load the function pointer directly from the struct field (single load)
	fnVal := cls.LoadField(bh, idx, fieldType)
	if fnVal == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("function pointer is nil for %s.%s", cls.Name, m.Member))
	}

	var funcType *types.FuncType
	if ptrType, ok := fieldType.(*types.PointerType); ok {
		funcType, ok = ptrType.ElemType.(*types.FuncType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("expected pointer-to-function, got pointer to %T", ptrType.ElemType))
		}
	} else {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("expected pointer-to-function type for field, got %T", fieldType))
	}

	// Build args
	args := make([]value.Value, 0, len(ex.Arguments)+1)
	actualArgCount := int64(len(ex.Arguments))

	// pass thread function
	args = append(args, fnVal)

	// pass argument count
	args = append(args, constant.NewInt(types.I32, actualArgCount))

	// pass rest of the args
	for i, argExp := range ex.Arguments[1:] {
		v := t.ProcessExpression(bh, argExp)
		raw := v.Load(bh)
		expected := funcType.Params[i]
		raw = t.st.TypeHandler.ImplicitTypeCast(bh, utils.GetTypeString(expected), raw)
		args = append(args, raw)
	}

	// pass `this` pointer as last arg
	thisPtr := cls.Load(bh)
	args = append(args, thisPtr)

	// Final call
	bh.N.NewCall(meth, args...)
	return nil

}

func (t *ExpressionHandler) callClassMethod(bh *bc.BlockHolder, ex ast.CallExpression, m ast.MemberExpression) tf.Var {
	// evaluate the base expression
	baseVar := t.ProcessExpression(bh, m.Member)
	if baseVar == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "nil base for member expression")
	}

	// validate baseVar
	cls, ok := baseVar.(*tf.Class)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "member access base is not a Class type")
	}
	if cls == nil || cls.Ptr == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "class or class.Ptr is nil for class")
	}

	// validate class registration
	classMeta := t.st.Classes[cls.Name]
	if classMeta == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "unknown class metadata: "+cls.Name)
	}

	// validate accessebility
	methodFqName := fmt.Sprintf("%s.%s", cls.Name, m.Property)
	idx, ok := classMeta.FieldIndexMap[methodFqName]

	// to check whether access is comming from method in own class or elsewhere.
	// this decides access scope of that function.
	if resolveRootMember(m) != constants.THIS {
		if _, ok := classMeta.InternalFields[methodFqName]; ok {
			errorutils.Abort(errorutils.FieldNotAccessible, cls.Name, m.Property)
		}
	}
	if !ok {
		errorutils.Abort(errorutils.UnknownMethod, m.Property)
	}

	fieldType := classMeta.StructType().Fields[idx]

	// Load the function pointer directly from the struct field (single load)
	fnVal := cls.LoadField(bh, idx, fieldType)
	if fnVal == nil {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("function pointer is nil for %s.%s", cls.Name, m.Member))
	}

	var funcType *types.FuncType
	if ptrType, ok := fieldType.(*types.PointerType); ok {
		funcType, ok = ptrType.ElemType.(*types.FuncType)
		if !ok {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("expected pointer-to-function, got pointer to %T", ptrType.ElemType))
		}
	} else {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, fmt.Sprintf("expected pointer-to-function type for field, got %T", fieldType))
	}

	// Build args
	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v := t.ProcessExpression(bh, argExp)
		raw := v.Load(bh)
		expected := classMeta.MethodArgs[methodFqName][i]
		raw = t.st.TypeHandler.ImplicitTypeCast(bh, expected.Get(), raw)
		args = append(args, raw)
	}

	// Append `this` pointer as last arg
	thisPtr := cls.Load(bh)
	args = append(args, thisPtr)

	// Call
	ret := bh.N.NewCall(fnVal, args...)

	// Return handling
	retType := funcType.RetType
	if retType == types.Void {
		return nil
	}

	// @todo: not tested
	tp := classMeta.Returns[methodFqName]
	return t.st.TypeHandler.BuildVar(bh, tf.NewType(tp.Get(), tp.GetUnderlyingType()), ret)

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
