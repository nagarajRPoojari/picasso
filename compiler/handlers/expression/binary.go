package expression

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	"github.com/nagarajRPoojari/x-lang/lexer"
)

type BinaryOperation func(block *ir.Block, l, r value.Value) (value.Value, error)

var arithmatic map[lexer.TokenKind]BinaryOperation
var comparision map[lexer.TokenKind]BinaryOperation
var logical map[lexer.TokenKind]BinaryOperation

// ProcessBinaryExpression handles evaluation of a binary operation expression.
//
// It recursively processes the left and right operands, performs type-safe casting,
// applies the corresponding operator (arithmetic, comparison, or logical),
// and produces a new runtime variable with the result.
//
// Supported categories:
//   - Arithmetic (e.g., +, -, *, /) → result as Float64
//   - Comparison (e.g., <, >, ==)   → result as Boolean (I1)
//   - Logical (e.g., &&, ||)        → result as Boolean (I1)
//
// Parameters:
//
//	block - the current IR block in which code generation should occur
//	ex    - the binary expression AST node
//
// Returns:
//
//	tf.Var     - the resulting variable (float64 or boolean depending on operator)
//	*ir.Block  - the (possibly updated) IR block after processing
func (t *ExpressionHandler) ProcessBinaryExpression(bh *bc.BlockHolder, ex ast.BinaryExpression) tf.Var {
	left := t.ProcessExpression(bh, ex.Left)

	right := t.ProcessExpression(bh, ex.Right)

	if left == nil || right == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	lv := left.Load(bh)
	rv := right.Load(bh)

	if op, ok := arithmatic[ex.Operator.Kind]; ok {
		f := &floats.Float64{}
		lvf, err := f.Cast(bh, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv.Type(), tf.FLOAT64)
		}
		rvf, err := f.Cast(bh, rv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, rv.Type(), tf.FLOAT64)
		}

		res, err := op(bh.N, lvf, rvf)

		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		ptr := bh.V.NewAlloca(types.Double)
		bh.N.NewStore(res, ptr)
		return &floats.Float64{NativeType: types.Double, Value: ptr}

	} else if op, ok := comparision[ex.Operator.Kind]; ok {
		// comparision operations are done on float type operands only
		f := &floats.Float64{}
		lvf, err := f.Cast(bh, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv.Type(), tf.FLOAT64)
		}
		rvf, err := f.Cast(bh, rv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, rv.Type(), tf.FLOAT64)
		}

		res, err := op(bh.N, lvf, rvf)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		ptr := bh.V.NewAlloca(types.I1)
		bh.N.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}

	} else if op, ok := logical[ex.Operator.Kind]; ok {
		f := &boolean.Boolean{}
		lvf, err := f.Cast(bh, lv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, lv.Type(), tf.BOOLEAN)
		}
		rvf, err := f.Cast(bh, rv)
		if err != nil {
			errorutils.Abort(errorutils.ImplicitTypeCastError, rv.Type(), tf.BOOLEAN)
		}

		res, err := op(bh.N, lvf, rvf)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		ptr := bh.V.NewAlloca(types.I1)
		bh.N.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}
	}

	errorutils.Abort(errorutils.InvalidBinaryExpressionOperator, ex.Operator.Value)
	return nil
}

// initOpLookUpTables inits lookup table mapping operand token with its
// corresponding operation
func initOpLookUpTables() {
	arithmatic = make(map[lexer.TokenKind]BinaryOperation)
	comparision = make(map[lexer.TokenKind]BinaryOperation)
	logical = make(map[lexer.TokenKind]BinaryOperation)

	arithmatic[lexer.PLUS] = add
	arithmatic[lexer.DASH] = sub
	arithmatic[lexer.STAR] = mul
	arithmatic[lexer.SLASH] = div

	comparision[lexer.LESS] = lt
	comparision[lexer.LESS_EQUALS] = lte
	comparision[lexer.GREATER] = gt
	comparision[lexer.GREATER_EQUALS] = gte
	comparision[lexer.EQUALS] = eq
	comparision[lexer.NOT_EQUALS] = ne

	logical[lexer.AND] = and
	logical[lexer.OR] = or
}

func add(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFAdd(lvf, rvf), nil
}
func sub(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFSub(lvf, rvf), nil
}
func mul(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFMul(lvf, rvf), nil
}
func div(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFDiv(lvf, rvf), nil
}

func eq(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredOEQ, lvf, rvf), nil
}
func lt(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredOLT, lvf, rvf), nil
}
func lte(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredOLE, lvf, rvf), nil
}
func gt(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredOGT, lvf, rvf), nil
}
func gte(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredOGE, lvf, rvf), nil
}
func ne(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewFCmp(enum.FPredONE, lvf, rvf), nil
}

func and(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewAnd(lvf, rvf), nil
}
func or(block *ir.Block, lvf, rvf value.Value) (value.Value, error) {
	return block.NewOr(lvf, rvf), nil
}
