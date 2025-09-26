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
		fn, ok := classMeta.Methods[methodKey]
		if !ok || fn == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: unknown method %s on class %s", m.Property, cls.Name))
		}

		// Build args for the user parameters (do not append `this` yet)
		args := make([]value.Value, 0, len(ex.Arguments)+1)
		for i, argExp := range ex.Arguments {
			v, safe := t.ProcessExpression(block, argExp)
			block = safe

			if v == nil {
				errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: nil arg %d for %s.%s", i, cls.Name, m.Property))
			}
			raw := v.Load(block)
			if raw == nil {
				errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: loaded nil arg %d for %s.%s", i, cls.Name, m.Property))
			}

			// If the callee expects a certain param type, cast to it
			expected := fn.Sig.Params[i]
			target := utils.GetTypeString(expected)
			raw, block = t.st.TypeHandler.ImplicitTypeCast(block, target, raw)
			if raw == nil {
				errorsx.PanicCompilationError(fmt.Sprintf(
					"handleCallExpression: ImplicitTypeCast returned nil for arg %d -> %s", i, target))
			}
			args = append(args, raw)

		}

		// Pass `this` as a pointer-to-struct (Slot returns pointer)
		thisPtr := cls.Slot()
		if thisPtr == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: this pointer is nil for %s", cls.Name))
		}

		// Check function expected param count: we declared 'this' last when creating fn,
		// adjust order according to how the function was declared.
		args = append(args, thisPtr)

		// ensure function is non-nil
		if fn == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: function pointer nil for %s.%s", cls.Name, m.Property))
		}

		// make function call
		ret := block.NewCall(fn, args...)
		if fn.Sig.RetType == types.Void {
			return nil, block
		}

		tp := utils.GetTypeString(fn.Sig.RetType)
		return t.st.TypeHandler.BuildVar(block, tf.Type(tp), ret), block
	}
	return nil, block
}
