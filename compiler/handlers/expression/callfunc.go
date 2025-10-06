package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
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
func (t *ExpressionHandler) CallFunc(block *ir.Block, ex ast.CallExpression) (tf.Var, *ir.Block) {
	// check if imported modules
	if m, ok := ex.Method.(ast.MemberExpression); ok {
		x, ok := m.Member.(ast.SymbolExpression)
		if ok {
			fName := fmt.Sprintf("%s.%s", x.Value, m.Property)
			if f, ok := t.st.LibMethods[fName]; ok {
				args := make([]tf.Var, 0)
				for _, v := range ex.Arguments {
					res, safe := t.ProcessExpression(block, v)
					block = safe
					args = append(args, res)
				}
				ret, safe := f(t.st.TypeHandler, t.st.Module, block, args)
				block = safe
				return ret, block
			}
		}

	}

	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		errorutils.Abort(errorutils.MemberExpressionError, "method call should be on instance")

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar, safe := t.ProcessExpression(block, m.Member)
		block = safe

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
		fnVal := cls.LoadField(block, idx, fieldType)
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
			v, safe := t.ProcessExpression(block, argExp)
			block = safe
			raw := v.Load(block)
			expected := funcType.Params[i]
			raw, block = t.st.TypeHandler.ImplicitTypeCast(block, utils.GetTypeString(expected), raw)
			args = append(args, raw)
		}

		// Append `this` pointer as last arg
		thisPtr := cls.Slot()
		args = append(args, thisPtr)

		// Call
		ret := block.NewCall(fnVal, args...)

		fmt.Printf("funcType.RetType: %v\n", funcType.RetType)

		// Return handling
		retType := funcType.RetType
		if retType == types.Void {
			return nil, block
		}

		tp := utils.GetTypeString(retType)
		var underlyingType string
		if tp == constants.ARRAY {
			tp := classMeta.Returns[methodKey]
			underlyingType = tp.(ast.ListType).GetEleType()
		}

		return t.st.TypeHandler.BuildVar(block, tf.Type(tp), ret, underlyingType), block

	}
	return nil, block
}
