package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// CallFunc handles invocation of a function or method call expression.
//
// Supports two main cases:
//  1. Calls to imported library methods (resolved via t.st.LibMethods)
//  2. Calls to class instance methods (resolved through class metadata and vtable-like lookup)
//
// Steps:
//   - If the call looks like a library call (e.g., module.func), resolve and invoke directly.
//   - Otherwise, evaluate the base expression, resolve the class metadata, locate the method
//     function pointer, perform argument type-casting, and emit a call instruction.
//   - Appends the `this` pointer automatically as the final argument for instance methods.
//   - Returns a wrapped runtime variable if the function has a return type; otherwise returns nil.
//
// Parameters:
//
//	block - the current IR block in which code generation should occur
//	ex    - the call expression AST node
//
// Returns:
//
//	tf.Var     - the resulting variable if the function returns a value, otherwise nil
//	*ir.Block  - the (possibly updated) IR block after processing
func (t *ExpressionHandler) CallFunc(bh *bc.BlockHolder, ex ast.CallExpression) tf.Var {
	// check if imported modules
	if m, ok := ex.Method.(ast.MemberExpression); ok {
		x, ok := m.Member.(ast.SymbolExpression)
		if ok {
			fName := fmt.Sprintf("%s.%s", x.Value, m.Property)
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
		// errorutils.Abort(errorutils.MemberExpressionError, "method call should be on instance")
		if m.Value == c.FUNC_THREAD {
			if meth, ok := t.st.CI.Funcs[m.Value]; ok {
				v := t.ProcessExpression(bh, ex.Arguments[0])
				raw := v.Load(bh)
				this := t.ProcessExpression(bh, ex.Arguments[0].(ast.MemberExpression).Member)
				bh.N.NewCall(meth, raw, this.Load(bh))
			}
			return nil
		}

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

		methodKey := t.st.IdentifierBuilder.Attach(cls.Name, m.Property)
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
