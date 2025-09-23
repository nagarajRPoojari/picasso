package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) ProcessPrefixExpression(block *ir.Block, ex ast.PrefixExpression) tf.Var {
	operand := t.ProcessExpression(block, ex.Operand)

	var res value.Value
	lv := operand.Load(block)

	switch ex.Operator.Value {
	case "-":
		f := &floats.Float64{}
		val, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to cast %s to float", lv))
		}
		res = block.NewFNeg(val)
	case "!":
		f := &boolean.Boolean{}
		val, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to cast %s to float", lv))
		}
		one := constant.NewInt(types.I1, 1)
		res = block.NewXor(val, one)
	}

	switch res.Type().Equal(types.I1) {
	case true:
		ptr := block.NewAlloca(types.I1)
		block.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}
	default:
		ptr := block.NewAlloca(types.Double)
		block.NewStore(res, ptr)
		return &floats.Float64{NativeType: types.Double, Value: ptr}
	}
}
