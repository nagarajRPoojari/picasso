package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) ProcessNewExpression(block *ir.Block, ex ast.NewExpression) tf.Var {
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
		x := t.ProcessExpression(block, exp.AssignedValue)

		fieldType := structType.Fields[index]
		instance.UpdateField(block, index, x.Load(block), fieldType)
	}

	t.CallConstructor(block, instance, ex.Instantiation)
	return instance
}
