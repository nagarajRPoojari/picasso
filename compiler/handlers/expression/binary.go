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
	"github.com/nagarajRPoojari/x-lang/lexer"
)

type Operation func(block *ir.Block, l, r value.Value) (value.Value, error)

var arithmatic map[lexer.TokenKind]Operation
var comparision map[lexer.TokenKind]Operation
var logical map[lexer.TokenKind]Operation

func initOpLookUpTables() {
	arithmatic = make(map[lexer.TokenKind]Operation)
	comparision = make(map[lexer.TokenKind]Operation)
	logical = make(map[lexer.TokenKind]Operation)

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

func (t *ExpressionHandler) ProcessBinaryExpression(block *ir.Block, ex ast.BinaryExpression) tf.Var {
	left := t.ProcessExpression(block, ex.Left)
	right := t.ProcessExpression(block, ex.Right)
	if left == nil || right == nil {
		errorsx.PanicCompilationError("nil operand in binary expression")
	}

	lv := left.Load(block)
	rv := right.Load(block)

	if op, ok := arithmatic[ex.Operator.Kind]; ok {
		f := &floats.Float64{}
		lvf, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		rvf, err := f.Cast(block, rv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}

		res, err := op(block, lvf, rvf)

		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		ptr := block.NewAlloca(types.Double)
		block.NewStore(res, ptr)
		return &floats.Float64{NativeType: types.Double, Value: ptr}

	} else if op, ok := comparision[ex.Operator.Kind]; ok {
		f := &floats.Float64{}
		lvf, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		rvf, err := f.Cast(block, rv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}

		res, err := op(block, lvf, rvf)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		ptr := block.NewAlloca(types.I1)
		block.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}

	} else if op, ok := logical[ex.Operator.Kind]; ok {
		f := &boolean.Boolean{}
		lvf, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		rvf, err := f.Cast(block, rv)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}

		res, err := op(block, lvf, rvf)
		if err != nil {
			errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d: %v", int(ex.Operator.Kind), err))
		}
		ptr := block.NewAlloca(types.I1)
		block.NewStore(res, ptr)
		return &boolean.Boolean{NativeType: types.I1, Value: ptr}
	}
	errorsx.PanicCompilationError(fmt.Sprintf("failed to do operation %d", int(ex.Operator.Kind)))
	return nil
}
