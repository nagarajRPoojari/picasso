package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
)

// ProcessPrefixExpression evaluates a unary (prefix) expression and returns the result.
//
// Supported operators:
//   - "-" : negates a numeric operand (Float64)
//   - "!" : logical NOT on a boolean operand (Boolean)
//
// Parameters:
//
//	block - current IR block
//	ex    - AST PrefixExpression node
//
// Returns:
//
//	tf.Var - result variable
func (t *ExpressionHandler) ProcessPrefixExpression(bh tf.BlockHolder, ex ast.PrefixExpression) tf.Var {
	operand, safe := t.ProcessExpression(bh, ex.Operand)
	bh = safe

	var res value.Value
	lv := operand.Load(bh.N)

	switch ex.Operator.Value {
	case "-":
		f := &floats.Float64{}
		val, err := f.Cast(bh.N, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv, tf.FLOAT64)
		}
		res = bh.N.NewFNeg(val)
	case "!":
		f := &boolean.Boolean{}
		val, err := f.Cast(bh.N, lv)
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
