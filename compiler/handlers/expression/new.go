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

func (t *ExpressionHandler) callConstructor(block *ir.Block, cls *tf.Class, ex ast.CallExpression) *ir.Block {
	// Get the method symbol and metadata
	m := ex.Method.(ast.SymbolExpression)
	meth := t.st.IdentifierBuilder.Attach(m.Value, m.Value)
	meta := t.st.Classes[cls.Name]
	idx := meta.FieldIndexMap[meth]

	// Get struct and field type
	st := meta.StructType()
	fieldType := st.Fields[idx]

	// Load the function pointer directly from the struct field (single load)
	fnVal := cls.LoadField(block, idx, fieldType)
	if fnVal == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: function pointer is nil for %s.%s", cls.Name, m.Value))
	}

	// Ensure field type is pointer-to-function
	var funcType *types.FuncType
	if ptrType, ok := fieldType.(*types.PointerType); ok {
		funcType, ok = ptrType.ElemType.(*types.FuncType)
		if !ok {
			panic(fmt.Sprintf("expected pointer-to-function, got pointer to %T", ptrType.ElemType))
		}
	} else {
		panic(fmt.Sprintf("expected pointer-to-function type for field, got %T", fieldType))
	}

	// Build arguments
	args := make([]value.Value, 0, len(ex.Arguments)+1)
	for i, argExp := range ex.Arguments {
		v, safe := t.ProcessExpression(block, argExp)
		block = safe
		if v == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: nil arg %d for %s", i, m.Value))
		}
		raw := v.Load(block)
		if raw == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: loaded nil arg %d for %s", i, m.Value))
		}

		// Implicit type cast if needed
		expected := funcType.Params[i]
		target := utils.GetTypeString(expected)
		raw, safe = t.st.TypeHandler.ImplicitTypeCast(block, target, raw)
		block = safe
		if raw == nil {
			errorsx.PanicCompilationError(fmt.Sprintf(
				"handleConstructorCall: ImplicitTypeCast returned nil for arg %d -> %s", i, target))
		}
		args = append(args, raw)
	}

	// Append `this` pointer as the last argument
	thisPtr := cls.Slot()
	if thisPtr == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("handleConstructorCall: this pointer is nil for %s", cls.Name))
	}
	args = append(args, thisPtr)

	// Call the function pointer
	block.NewCall(fnVal, args...)

	return block
}

func (t *ExpressionHandler) ProcessNewExpression(block *ir.Block, ex ast.NewExpression) (tf.Var, *ir.Block) {
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

	for name, index := range meta.FieldIndexMap {
		exp, ok := meta.VarAST[name]
		if !ok {
			f := t.st.Classes[meth.Value].Methods[name]
			fieldType := t.st.Classes[meth.Value].StructType().Fields[index]
			instance.UpdateField(block, index, f, fieldType)
			continue
		}
		fieldType := structType.Fields[index]

		var v tf.Var
		if exp.AssignedValue == nil {
			v = t.st.TypeHandler.BuildVar(block, tf.Type(exp.ExplicitType.Get()), nil)
		} else {
			_v, safe := t.ProcessExpression(block, exp.AssignedValue)
			v = _v
			block = safe

			casted, safe := t.st.TypeHandler.ImplicitTypeCast(block, exp.ExplicitType.Get(), v.Load(block))
			block = safe

			v = t.st.TypeHandler.BuildVar(block, tf.Type(exp.ExplicitType.Get()), casted)
		}
		instance.UpdateField(block, index, v.Load(block), fieldType)
		t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	block = t.callConstructor(block, instance, ex.Instantiation)
	return instance, block
}
