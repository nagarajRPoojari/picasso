package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
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
	if m, ok := ex.Method.(ast.MemberExpression); ok {
		x, ok := m.Member.(ast.SymbolExpression)
		if ok {
			// base module func calls are strictly expected to be in module_name.func() format. e.g, syncio.printf
			fName := fmt.Sprintf("%s.%s", t.st.Imports[x.Value].Name, m.Property)
			if f, ok := t.st.LibMethods[fName]; ok {
				args := make([]tf.Var, 0)
				for _, v := range ex.Arguments {
					res := t.ProcessExpression(bh, v)
					args = append(args, res)
				}
				ret := f(t.st.TypeHandler, t.st.Module, bh, args)
				return ret
			}
		}
	}

	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		// @fix: refactor this piece of code
		// thread() special functions comes from base modules but exception for module_name.func() format
		// in future i might add many such special funcs, so need to be kept somewhere else.
		if m.Value == c.FUNC_THREAD {
			if meth, ok := t.st.CI.Funcs[m.Value]; ok {
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

				methodKey := fmt.Sprintf("%s.%s", cls.Name, m.Property)
				idx, ok := classMeta.FieldIndexMap[methodKey]

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
				args = append(args, t.ProcessExpression(bh, ex.Arguments[0]).Load(bh))

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
		}

		// support for direct c call.
		// @fix: need to be removed as it is completely unsafe.
		if meth, ok := t.st.CI.Funcs[m.Value]; ok {
			args := make([]value.Value, 0, len(ex.Arguments)+1)
			for _, argExp := range ex.Arguments {
				v := t.ProcessExpression(bh, argExp)
				raw := v.Load(bh)
				args = append(args, raw)
			}
			ret := bh.N.NewCall(meth, args...)
			return t.st.TypeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(ret.Type())), ret)
		}

		errorutils.Abort(errorutils.UnknownMethod, m.Value)

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar := t.ProcessExpression(bh, m.Member)

		if baseVar == nil {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalFuncCallError, "nil base for member expression")
		}

		cls, ok := baseVar.(*tf.Class)
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

		methodKey := fmt.Sprintf("%s.%s", cls.Name, m.Property)
		idx, ok := classMeta.FieldIndexMap[methodKey]

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
		for i, argExp := range ex.Arguments {
			v := t.ProcessExpression(bh, argExp)
			raw := v.Load(bh)
			expected := funcType.Params[i]
			raw = t.st.TypeHandler.ImplicitTypeCast(bh, utils.GetTypeString(expected), raw)
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
		tp := classMeta.Returns[methodKey]
		return t.st.TypeHandler.BuildVar(bh, tf.NewType(tp.Get(), tp.GetUnderlyingType()), ret)

	}
	return nil
}
