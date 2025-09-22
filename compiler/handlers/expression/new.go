package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

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
		}
		instance.UpdateField(block, index, v.Load(block), fieldType)
		t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	t.CallConstructor(block, instance, ex.Instantiation)
	return instance
}
