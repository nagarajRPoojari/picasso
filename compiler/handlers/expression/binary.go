package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *ExpressionHandler) ProcessBinaryExpression(block *ir.Block, ex ast.BinaryExpression) tf.Var {
	left := t.ProcessExpression(block, ex.Left)
	right := t.ProcessExpression(block, ex.Right)
	if left == nil || right == nil {
		errorsx.PanicCompilationError("nil operand in binary expression")
	}

	lv := left.Load(block)
	rv := right.Load(block)

	// For arithmetic and comparison, cast to float
	f := &floats.Float64{}
	lvf, err := f.Cast(block, lv)
	if err != nil {
		errorsx.PanicCompilationError(fmt.Sprintf("failed to cast %s to float", lv))
	}
	rvf, err := f.Cast(block, rv)
	if err != nil {
		errorsx.PanicCompilationError(fmt.Sprintf("failed to cast %s to float", rv))
	}

	var res value.Value
	switch ex.Operator.Value {
	// Arithmetic
	case "+":
		res = block.NewFAdd(lvf, rvf)
	case "-":
		res = block.NewFSub(lvf, rvf)
	case "*":
		res = block.NewFMul(lvf, rvf)
	case "/":
		res = block.NewFDiv(lvf, rvf)

	// Comparisons (return i1 bools)
	case "==":
		res = block.NewFCmp(enum.FPredOEQ, lvf, rvf)
	case "!=":
		res = block.NewFCmp(enum.FPredONE, lvf, rvf)
	case "<":
		res = block.NewFCmp(enum.FPredOLT, lvf, rvf)
	case "<=":
		res = block.NewFCmp(enum.FPredOLE, lvf, rvf)
	case ">":
		res = block.NewFCmp(enum.FPredOGT, lvf, rvf)
	case ">=":
		res = block.NewFCmp(enum.FPredOGE, lvf, rvf)

	// Boolean (assume operands already i1)
	case "&&":
		res = block.NewAnd(lv, rv)
	case "||":
		res = block.NewOr(lv, rv)

	default:
		panic(fmt.Sprintf("unsupported operator: %s", ex.Operator.Value))
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
