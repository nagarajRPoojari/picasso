package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) CallConstructor(block *ir.Block, cls *tf.Class, ex ast.CallExpression) {
	m := ex.Method.(ast.SymbolExpression)
	meth := t.st.IdentifierBuilder.Attach(m.Value, m.Value)
	fn := t.st.Methods[meth]

	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v := t.ProcessExpression(block, argExp)
		if v == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: nil arg %d for %s", i, m.Value))
		}
		raw := v.Load(block)
		if raw == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: loaded nil arg %d for %s", i, m.Value))
		}

		expected := fn.Sig.Params[i]
		target := utils.GetTypeString(expected)
		raw = t.st.TypeHandler.CastToType(block, target, raw)
		if raw == nil {
			errorsx.PanicCompilationError(fmt.Sprintf(
				"handleCallExpression: CastToType returned nil for arg %d -> %s", i, target))
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

	block.NewCall(fn, args...)
}
