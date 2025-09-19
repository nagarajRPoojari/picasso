package floats

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Float16/Half stores 2 byte floating point value
type Float16 struct {
	NativeType *types.FloatType
	Value      value.Value
	GoVal      float32
}

func NewFloat16Var(block *ir.Block, init float32) *Float16 {
	slot := block.NewAlloca(types.Half)
	block.NewStore(constant.NewFloat(types.Half, float64(init)), slot)
	return &Float16{NativeType: types.Half, Value: slot, GoVal: init}
}

func (f *Float16) Update(block *ir.Block, v value.Value) { block.NewStore(v, f.Value) }
func (f *Float16) Load(block *ir.Block) value.Value      { return block.NewLoad(types.Half, f.Value) }
func (f *Float16) Constant() constant.Constant {
	return constant.NewFloat(types.Half, float64(f.GoVal))
}
func (f *Float16) Slot() value.Value { return f.Value }
func (f *Float16) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Half), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.NewFPTrunc(v, types.Half), nil
		case types.Half:
			return block.NewFPExt(v, types.Half), nil
		case types.Float:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
}
func (c *Float16) Type() types.Type         { return c.NativeType }
func (f *Float16) NativeTypeString() string { return "float16" }
