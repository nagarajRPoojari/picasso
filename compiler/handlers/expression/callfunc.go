package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

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
				return f(t.st.TypeHandler, t.st.Module, block, args), block
			}
		}

	}

	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		errorsx.PanicCompilationError("method call should be on instance")

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar, safe := t.ProcessExpression(block, m.Member)
		block = safe

		if baseVar == nil {
			errorsx.PanicCompilationError("handleCallExpression: nil baseVar for member expression")
		}

		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: member access base is not Class (got %T)", baseVar))
		}
		if cls == nil || cls.Ptr == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: class or class.Ptr is nil for class %v", cls))
		}

		classMeta := t.st.Classes[cls.Name]
		if classMeta == nil {
			errorsx.PanicCompilationError("handleCallExpression: unknown class metadata: " + cls.Name)
		}

		methodKey := t.st.IdentifierBuilder.Attach(cls.Name, m.Property)
		idx, ok := classMeta.FieldIndexMap[methodKey]

		if !ok {
			panic(fmt.Sprintf("unable to find method: %s", m.Property))
		}

		st := classMeta.StructType()
		fieldType := st.Fields[idx]

		// Load the function pointer directly from the struct field (single load)
		fnVal := cls.LoadField(block, idx, fieldType)
		if fnVal == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: function pointer is nil for %s.%s", cls.Name, m.Member))
		}

		var funcType *types.FuncType
		if ptrType, ok := fieldType.(*types.PointerType); ok {
			funcType, ok = ptrType.ElemType.(*types.FuncType)
			if !ok {
				panic(fmt.Sprintf("expected pointer-to-function, got pointer to %T", ptrType.ElemType))
			}
		} else {
			panic(fmt.Sprintf("expected pointer-to-function type for field, got %T", fieldType))
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

		// Return handling
		retType := funcType.RetType
		if retType == types.Void {
			return nil, block
		}

		tp := utils.GetTypeString(retType)
		return t.st.TypeHandler.BuildVar(block, tf.Type(tp), ret), block

	}
	return nil, block
}
