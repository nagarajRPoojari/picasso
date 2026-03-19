package expression

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/boolean"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/floats"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/ints"
	"github.com/nagarajRPoojari/picasso/irgen/lexer"
)

type BinaryOperation func(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error)

var arithmatic map[lexer.TokenKind]BinaryOperation
var comparision map[lexer.TokenKind]BinaryOperation
var logical map[lexer.TokenKind]BinaryOperation
var bitwise map[lexer.TokenKind]BinaryOperation

type ArithKind int

const (
	KindInvalid ArithKind = iota
	KindSignedInt
	KindUnsignedInt
	KindFloat
	KindPointer
)

func classifyVar(v tf.Var) ArithKind {
	switch v.(type) {

	// unsigned
	case *ints.UInt8, *ints.UInt16, *ints.UInt32, *ints.UInt64:
		return KindUnsignedInt

	// signed
	case *ints.Int8, *ints.Int16, *ints.Int32, *ints.Int64:
		return KindSignedInt

	// float
	case *floats.Float16, *floats.Float32, *floats.Float64:
		return KindFloat

	// pointer
	case *tf.Array, *tf.Class, *tf.String, *tf.NullVar, *tf.InterfaceH:
		return KindPointer
	}

	return KindInvalid
}

func commonKind(a, b ArithKind) ArithKind {
	// Pointer rules: pointers only with pointers
	if a == KindPointer || b == KindPointer {
		if a == KindPointer && b == KindPointer {
			return KindPointer
		}
		return KindInvalid
	}

	// Unsigned integers: only compatible with other unsigned integers
	if a == KindUnsignedInt || b == KindUnsignedInt {
		if a == KindUnsignedInt && b == KindUnsignedInt {
			return KindUnsignedInt
		}
		return KindInvalid
	}

	// Float dominates only signed integers
	if a == KindFloat || b == KindFloat {
		if a == KindSignedInt || b == KindSignedInt || a == KindFloat || b == KindFloat {
			return KindFloat
		}
		return KindInvalid
	}

	// Signed integers only with signed integers
	if a == KindSignedInt && b == KindSignedInt {
		return KindSignedInt
	}

	return KindInvalid
}

func normalizeOperands(th *tf.TypeHandler, bh *bc.BlockHolder, lv, rv tf.Var) (value.Value, value.Value, ArithKind, error) {

	lk := classifyVar(lv)
	rk := classifyVar(rv)

	k := commonKind(lk, rk)
	if k == KindInvalid {
		fmt.Printf("lk: %v\n", lk)
		fmt.Printf("rk: %v\n", rk)
		return nil, nil, KindInvalid,
			fmt.Errorf("incompatible operands for operation")
	}

	l := lv.Load(bh)
	r := rv.Load(bh)

	switch k {
	case KindFloat:
		return th.ImplicitFloatCast(bh, l, types.Double),
			th.ImplicitFloatCast(bh, r, types.Double),
			k, nil

	case KindUnsignedInt:
		return th.ImplicitUnsignedIntCast(bh, l, types.I64),
			th.ImplicitUnsignedIntCast(bh, r, types.I64),
			k, nil

	case KindSignedInt:
		return th.ImplicitIntCast(bh, l, types.I64),
			th.ImplicitIntCast(bh, r, types.I64),
			k, nil

	case KindPointer:
		return l, r, k, nil
	}

	return nil, nil, KindInvalid, fmt.Errorf("unreachable")
}

func buildFloat64FromValue(bh *bc.BlockHolder, v value.Value) tf.Var {
	ptr := bh.V.NewAlloca(types.Double)
	bh.N.NewStore(v, ptr)
	return &floats.Float64{NativeType: types.Double, Value: ptr}
}

func buildUnsignedInt64FromValue(bh *bc.BlockHolder, v value.Value) tf.Var {
	ptr := bh.V.NewAlloca(types.I64)
	bh.N.NewStore(v, ptr)
	return &ints.UInt64{NativeType: types.I64, Value: ptr}
}

func buildSignedInt64FromValue(bh *bc.BlockHolder, v value.Value) tf.Var {
	ptr := bh.V.NewAlloca(types.I64)
	bh.N.NewStore(v, ptr)
	return &ints.Int64{NativeType: types.I64, Value: ptr}
}

func buildBooleanFromValue(bh *bc.BlockHolder, v value.Value) tf.Var {
	ptr := bh.V.NewAlloca(types.I1)
	bh.N.NewStore(v, ptr)
	return &boolean.Boolean{NativeType: types.I1, Value: ptr}
}

// ProcessBinaryExpression generates LLVM IR for operations involving two operands.
// It handles arithmetic, comparison, and logical operators by performing the
// necessary type promotions (e.g., coercing numeric types to float64) and
// emitting the corresponding LLVM instructions.
//
// Key Logic:
//   - Evaluation Order: Recursively processes Left and Right expressions before
//     performing the operation.
//   - Type Coercion: Standardizes numeric operations to Double (float64) and
//     logical operations to i1 (boolean) to ensure ABI compatibility.
//   - Memory Allocation: Automatically allocates stack space (alloca) for the
//     result, returning a wrapped tf.Var for subsequent use in the pipeline.
func (t *ExpressionHandler) ProcessBinaryExpression(bh *bc.BlockHolder, ex ast.BinaryExpression) tf.Var {
	if op, ok := arithmatic[ex.Operator.Kind]; ok {
		res, err := op(t.st.TypeHandler, bh, ex.Left, ex.Right)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		return res

	} else if op, ok := comparision[ex.Operator.Kind]; ok {
		res, err := op(t.st.TypeHandler, bh, ex.Left, ex.Right)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		return res

	} else if op, ok := logical[ex.Operator.Kind]; ok {
		res, err := op(t.st.TypeHandler, bh, ex.Left, ex.Right)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		return res
	} else if op, ok := bitwise[ex.Operator.Kind]; ok {
		res, err := op(t.st.TypeHandler, bh, ex.Left, ex.Right)
		if err != nil {
			errorutils.Abort(errorutils.BinaryOperationError, err.Error())
		}
		return res
	}

	errorutils.Abort(errorutils.InvalidBinaryExpressionOperator, ex.Operator.Value)
	return nil
}

// initOpLookUpTables inits lookup table mapping operand token with its
// corresponding operation
func initOpLookUpTables(ex *ExpressionHandler) {
	arithmatic = make(map[lexer.TokenKind]BinaryOperation)
	comparision = make(map[lexer.TokenKind]BinaryOperation)
	logical = make(map[lexer.TokenKind]BinaryOperation)
	bitwise = make(map[lexer.TokenKind]BinaryOperation)

	arithmatic[lexer.PLUS] = ex.add
	arithmatic[lexer.DASH] = ex.sub
	arithmatic[lexer.STAR] = ex.mul
	arithmatic[lexer.SLASH] = ex.div
	arithmatic[lexer.PERCENT] = ex.mod

	comparision[lexer.LESS] = ex.lt
	comparision[lexer.LESS_EQUALS] = ex.lte
	comparision[lexer.GREATER] = ex.gt
	comparision[lexer.GREATER_EQUALS] = ex.gte
	comparision[lexer.EQUALS] = ex.eq
	comparision[lexer.NOT_EQUALS] = ex.ne

	logical[lexer.AND] = ex.logicalAnd
	logical[lexer.OR] = ex.logicalOr
	logical[lexer.QUESTION] = ex.logicalInstanceOf

	bitwise[lexer.BITWISE_AND] = ex.bitwiseAND
	bitwise[lexer.BITWISE_OR] = ex.bitwiseOR
	bitwise[lexer.BITWISE_XOR] = ex.bitwiseXOR
	bitwise[lexer.BITWIZE_LEFTSHIFT] = ex.bitwiseLeftShift
	bitwise[lexer.BITWIZE_RIGHTSHIFT] = ex.bitwiseRightShift
}

func (t *ExpressionHandler) add(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		lf := th.ImplicitFloatCast(bh, l, types.Double)
		rf := th.ImplicitFloatCast(bh, r, types.Double)
		return buildFloat64FromValue(bh, bh.N.NewFAdd(lf, rf)), nil
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewAdd(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewAdd(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("addtion not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported add operands")
}

func (t *ExpressionHandler) sub(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		lf := th.ImplicitFloatCast(bh, l, types.Double)
		rf := th.ImplicitFloatCast(bh, r, types.Double)
		return buildFloat64FromValue(bh, bh.N.NewFSub(lf, rf)), nil
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewSub(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewSub(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("subtraction not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported sub operands")
}

func (t *ExpressionHandler) mul(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		lf := th.ImplicitFloatCast(bh, l, types.Double)
		rf := th.ImplicitFloatCast(bh, r, types.Double)
		return buildFloat64FromValue(bh, bh.N.NewFMul(lf, rf)), nil
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewMul(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewMul(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("multiplication not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported mul operands")
}

func (t *ExpressionHandler) div(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		lf := th.ImplicitFloatCast(bh, l, types.Double)
		rf := th.ImplicitFloatCast(bh, r, types.Double)
		return buildFloat64FromValue(bh, bh.N.NewFDiv(lf, rf)), nil
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewSDiv(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewUDiv(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("multiplication not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported mul operands")
}

func (t *ExpressionHandler) mod(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat, KindSignedInt, KindUnsignedInt:
		lf := th.ImplicitFloatCast(bh, l, types.Double)
		rf := th.ImplicitFloatCast(bh, r, types.Double)

		return buildFloat64FromValue(bh, bh.N.NewFRem(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("modulo not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported mod operands")
}

func (t *ExpressionHandler) eq(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredOEQ, l, r)), nil
	case KindSignedInt, KindUnsignedInt, KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredEQ, l, r)), nil
	}

	return nil, fmt.Errorf("unsupported eq")
}

func (t *ExpressionHandler) ne(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredONE, l, r)), nil
	case KindSignedInt, KindUnsignedInt, KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredNE, l, r)), nil
	}

	return nil, fmt.Errorf("unsupported ne")
}

func (t *ExpressionHandler) lt(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredOLT, l, r)), nil
	case KindSignedInt:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredSLT, l, r)), nil
	case KindUnsignedInt:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredULT, l, r)), nil
	case KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredULT, l, r)), nil
	}
	return nil, fmt.Errorf("unsupported lt")
}

func (t *ExpressionHandler) lte(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredOLE, l, r)), nil
	case KindSignedInt:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredSLE, l, r)), nil
	case KindUnsignedInt, KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredULE, l, r)), nil
	}
	return nil, fmt.Errorf("unsupported lte")
}

func (t *ExpressionHandler) gt(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredOGT, l, r)), nil
	case KindSignedInt:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredSGT, l, r)), nil
	case KindUnsignedInt, KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredUGT, l, r)), nil
	}
	return nil, fmt.Errorf("unsupported gt")
}

func (t *ExpressionHandler) gte(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return buildBooleanFromValue(bh, bh.N.NewFCmp(enum.FPredOGE, l, r)), nil
	case KindSignedInt:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredSGE, l, r)), nil
	case KindUnsignedInt, KindPointer:
		return buildBooleanFromValue(bh, bh.N.NewICmp(enum.IPredUGE, l, r)), nil
	}
	return nil, fmt.Errorf("unsupported gte")
}

func toBool(_ *tf.TypeHandler, bh *bc.BlockHolder, v tf.Var) (value.Value, error) {

	val := v.Load(bh)

	switch t := val.Type().(type) {

	case *types.IntType:
		return bh.N.NewICmp(
			enum.IPredNE,
			val,
			constant.NewInt(t, 0),
		), nil

	case *types.FloatType:
		zero := constant.NewFloat(t, 0.0)
		notZero := bh.N.NewFCmp(enum.FPredONE, val, zero)
		notNaN := bh.N.NewFCmp(enum.FPredORD, val, val)
		return bh.N.NewAnd(notZero, notNaN), nil

	case *types.PointerType:
		return bh.N.NewICmp(
			enum.IPredNE,
			val,
			constant.NewNull(t),
		), nil

	default:
		return nil, fmt.Errorf("cannot convert %T to bool", t)
	}
}

func (t *ExpressionHandler) logicalAnd(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	// Short-circuit evaluation: if left is false, result is false without evaluating right

	// Get the parent function to create new blocks
	fn := bh.N.Parent

	// Evaluate left operand
	lv := t.ProcessExpression(bh, lex)
	if lv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	lb, err := toBool(th, bh, lv)
	if err != nil {
		return nil, err
	}

	// Create blocks for short-circuit logic
	rhsBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	leftBlock := bh.N

	bh.N.NewCondBr(lb, rhsBlock.N, endBlock.N)
	rv := t.ProcessExpression(rhsBlock, rex)
	if rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	rb, err := toBool(th, rhsBlock, rv)
	if err != nil {
		return nil, err
	}
	rhsBlock.N.NewBr(endBlock.N)

	// Create phi node in end block to merge results
	// If we came from leftBlock (left was false), result is false
	// If we came from rhsBlock (left was true), result is rb
	phi := endBlock.N.NewPhi(ir.NewIncoming(constant.False, leftBlock), ir.NewIncoming(rb, rhsBlock.N))

	// Update block holder to continue from end block
	bh.Update(endBlock.V, endBlock.N)

	return buildBooleanFromValue(bh, phi), nil
}

func (t *ExpressionHandler) logicalOr(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	fn := bh.N.Parent

	lv := t.ProcessExpression(bh, lex)
	if lv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	lb, err := toBool(th, bh, lv)
	if err != nil {
		return nil, err
	}

	rhsBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	leftBlock := bh.N

	bh.N.NewCondBr(lb, endBlock.N, rhsBlock.N)

	rv := t.ProcessExpression(rhsBlock, rex)
	if rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	rb, err := toBool(th, rhsBlock, rv)
	if err != nil {
		return nil, err
	}

	rhsBlock.N.NewBr(endBlock.N)

	// Create phi node in end block to merge results
	// If we came from leftBlock (left was true), result is true
	// If we came from rhsBlock (left was false), result is rb
	phi := endBlock.N.NewPhi(ir.NewIncoming(constant.True, leftBlock), ir.NewIncoming(rb, rhsBlock.N))

	// Update block holder to continue from end block
	bh.Update(endBlock.V, endBlock.N)

	return buildBooleanFromValue(bh, phi), nil
}

func (t *ExpressionHandler) logicalInstanceOf(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	if lvc, ok := lv.(*tf.Class); ok {
		if rvc, ok := rv.(*tf.Class); ok {
			// both are class
			if lvc.Name == rvc.Name {
				return buildBooleanFromValue(bh, constant.True), nil
			}
		}
		if rvc, ok := rv.(*tf.InterfaceH); ok {
			// rhs is interface
			if th.ClassUDTS[lvc.Name].Implements == rvc.Name {
				return buildBooleanFromValue(bh, constant.True), nil
			}
		}
		return buildBooleanFromValue(bh, constant.False), nil
	}
	return nil, fmt.Errorf("invalid left hand operand for instance of")
}

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

	switch ex.Operator.Value {
	case "-":
		res, err := neg(t.st.TypeHandler, bh, operand)
		if err != nil {
			errorutils.Abort(errorutils.PrefixOperationError, err.Error())
		}
		return res
	case "!":
		res, err := not(t.st.TypeHandler, bh, operand)
		if err != nil {
			errorutils.Abort(errorutils.PrefixOperationError, err.Error())
		}
		return res
	case "~":
		res, err := bitwiseNOT(t.st.TypeHandler, bh, operand)
		if err != nil {
			errorutils.Abort(errorutils.PrefixOperationError, err.Error())
		}
		return res
	}

	panic(fmt.Sprintf("invalid prefix operation: %s", ex.Operator.Value))
}

func neg(th *tf.TypeHandler, bh *bc.BlockHolder, v tf.Var) (tf.Var, error) {
	tp := classifyVar(v)
	switch tp {
	case KindSignedInt:
		val := th.ImplicitFloatCast(bh, v.Load(bh), types.Double)
		return buildFloat64FromValue(bh, bh.N.NewFNeg(val)), nil
	case KindFloat:
		val := v.Load(bh)
		return buildFloat64FromValue(bh, bh.N.NewFNeg(val)), nil
	case KindUnsignedInt:
		return nil, fmt.Errorf("negation not allowed on unsigned dtypes")
	case KindPointer:
		return nil, fmt.Errorf("negation not allowed on unsigned dtypes")
	}

	return nil, fmt.Errorf("unsupported neg")
}

func not(th *tf.TypeHandler, bh *bc.BlockHolder, v tf.Var) (tf.Var, error) {
	vb, err := toBool(th, bh, v)
	if err != nil {
		return nil, err
	}
	one := constant.NewInt(types.I1, 1)
	return buildBooleanFromValue(bh, bh.N.NewXor(vb, one)), nil
}

func (t *ExpressionHandler) bitwiseOR(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return nil, fmt.Errorf("bitwise or not allowed on float types")
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewOr(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewOr(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("bitwise or not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported bitwise or operands")
}

func (t *ExpressionHandler) bitwiseXOR(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return nil, fmt.Errorf("bitwise XOR not allowed on float types")
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewXor(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewXor(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("bitwise xor not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported bitwise xor operands")
}

func (t *ExpressionHandler) bitwiseAND(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)

	rv := t.ProcessExpression(bh, rex)

	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}

	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}

	switch k {
	case KindFloat:
		return nil, fmt.Errorf("bitwise and not allowed on float types")
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewAnd(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewAnd(lf, rf)), nil

	case KindPointer:
		return nil, fmt.Errorf("bitwise and not allowed on pointer types")
	}

	return nil, fmt.Errorf("unsupported bitwise and operands")
}

func (t *ExpressionHandler) bitwiseLeftShift(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)
	rv := t.ProcessExpression(bh, rex)
	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}
	switch k {
	case KindFloat:
		return nil, fmt.Errorf("bitwise left shift not allowed on float types")
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewShl(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewShl(lf, rf)), nil
	case KindPointer:
		return nil, fmt.Errorf("bitwise left shift not allowed on pointer types")
	}
	return nil, fmt.Errorf("unsupported bitwise left shift operands")
}

func (t *ExpressionHandler) bitwiseRightShift(th *tf.TypeHandler, bh *bc.BlockHolder, lex, rex ast.Expression) (tf.Var, error) {
	lv := t.ProcessExpression(bh, lex)
	rv := t.ProcessExpression(bh, rex)
	if lv == nil || rv == nil {
		errorutils.Abort(errorutils.InvalidBinaryExpressionOperand)
	}
	l, r, k, err := normalizeOperands(th, bh, lv, rv)
	if err != nil {
		return nil, err
	}
	switch k {
	case KindFloat:
		return nil, fmt.Errorf("bitwise right shift not allowed on float types")
	case KindSignedInt:
		lf := th.ImplicitIntCast(bh, l, types.I64)
		rf := th.ImplicitIntCast(bh, r, types.I64)
		return buildSignedInt64FromValue(bh, bh.N.NewAShr(lf, rf)), nil
	case KindUnsignedInt:
		lf := th.ImplicitUnsignedIntCast(bh, l, types.I64)
		rf := th.ImplicitUnsignedIntCast(bh, r, types.I64)
		return buildUnsignedInt64FromValue(bh, bh.N.NewLShr(lf, rf)), nil
	case KindPointer:
		return nil, fmt.Errorf("bitwise right shift not allowed on pointer types")
	}
	return nil, fmt.Errorf("unsupported bitwise right shift operands")
}

func toInt64(_ *tf.TypeHandler, bh *bc.BlockHolder, v tf.Var) (value.Value, error) {
	val := v.Load(bh)
	i64 := types.I64

	switch t := val.Type().(type) {

	case *types.IntType:
		// Extend or truncate to i64
		if t.BitSize < 64 {
			return bh.N.NewSExt(val, i64), nil
		}
		if t.BitSize > 64 {
			return bh.N.NewTrunc(val, i64), nil
		}
		return val, nil

	case *types.FloatType:
		// float -> int64
		return bh.N.NewFPToSI(val, i64), nil

	case *types.PointerType:
		// pointer -> int64
		return bh.N.NewPtrToInt(val, i64), nil
	default:
		return nil, fmt.Errorf("cannot convert %T to int64", t)
	}
}

// @todo: test this
func bitwiseNOT(th *tf.TypeHandler, bh *bc.BlockHolder, v tf.Var) (tf.Var, error) {
	iv, err := toInt64(th, bh, v)
	if err != nil {
		return nil, err
	}

	allOnes := constant.NewInt(types.I64, -1)
	notv := bh.N.NewXor(iv, allOnes)

	return buildUnsignedInt64FromValue(bh, notv), nil
}
