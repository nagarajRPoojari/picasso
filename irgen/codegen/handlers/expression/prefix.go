package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/type/primitives/boolean"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/type/primitives/floats"
)

// ProcessPrefixExpression generates LLVM IR for unary operations such as
// numerical negation (-) and logical NOT (!). It evaluates the operand,
// performs the necessary type promotion, and applies the operator using
// specific LLVM instructions like FNeg or Xor.
//
// Technical Logic:
//   - Numerical Negation: Coerces the operand to a double-precision float
//     and emits an 'fneg' instruction to maintain consistency with the
//     compiler's floating-point-first arithmetic strategy.
//   - Logical NOT: Coerces the operand to a 1-bit integer (i1) and performs
//     an 'xor' operation with a constant 1 (true) to flip the boolean state.
//   - Result Storage: Allocates a new stack slot (alloca) for the result
//     to ensure the returned value is addressable as a tf.Var.
func (t *ExpressionHandler) ProcessPrefixExpression(bh *bc.BlockHolder, ex ast.PrefixExpression) tf.Var {
	operand := t.ProcessExpression(bh, ex.Operand)

	var res value.Value
	lv := operand.Load(bh)

	switch ex.Operator.Value {
	case "-":
		f := &floats.Float64{}
		val, err := f.Cast(bh, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv, tf.FLOAT64)
		}
		res = bh.N.NewFNeg(val)
	case "!":
		f := &boolean.Boolean{}
		val, err := f.Cast(bh, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv, tf.FLOAT64)
		}
		one := constant.NewInt(types.I1, 1)
		res = bh.N.NewXor(val, one)
	}

	switch res.Type().Equal(types.I1) {
	case true:
		ptr := bh.V.NewAlloca(types.I1)
		bh.N.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}
	default:
		ptr := bh.V.NewAlloca(types.Double)
		bh.N.NewStore(res, ptr)
		return &floats.Float64{NativeType: types.Double, Value: ptr}
	}
}
