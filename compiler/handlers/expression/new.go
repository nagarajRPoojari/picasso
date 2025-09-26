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

func (t *ExpressionHandler) callConstructor(block *ir.Block, cls *tf.Class, ex ast.CallExpression) {
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
		raw, safe := t.st.TypeHandler.ImplicitTypeCast(block, target, raw)
		block = safe
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

	block.NewCall(fn, args...)
}

func (t *ExpressionHandler) ProcessNewExpression(block *ir.Block, ex ast.NewExpression) tf.Var {
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	meth := ex.Instantiation.Method.(ast.SymbolExpression)
	classMeta := t.st.Classes[meth.Value]
	if classMeta == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("unknown class: %s", meth.Value))
	}

	instance := tf.NewClass(block, meth.Value, classMeta.UDT)
	structType := classMeta.StructType()
	meta := t.st.Classes[meth.Value]

	for name, index := range meta.VarIndexMap {
		exp := meta.VarAST[name]
		fieldType := structType.Fields[index]

		var v tf.Var
		if exp.AssignedValue == nil {
			v = t.st.TypeHandler.BuildVar(block, tf.Type(exp.ExplicitType.Get()), nil)
		} else {
			v = t.ProcessExpression(block, exp.AssignedValue)
			casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, exp.ExplicitType.Get(), v.Load(block))
			block = safe
			v = t.st.TypeHandler.BuildVar(block, tf.Type(exp.ExplicitType.Get()), casted)
		}
		instance.UpdateField(block, index, v.Load(block), fieldType)
		t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	t.callConstructor(block, instance, ex.Instantiation)
	return instance
}
