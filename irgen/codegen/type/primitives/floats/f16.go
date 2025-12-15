package floats

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"

	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
)

// Float16/Half stores 2 byte floating point value
type Float16 struct {
	NativeType *types.FloatType
	Value      value.Value
	GoVal      float32
}

func NewFloat16Var(block *bc.BlockHolder, init float32) *Float16 {
	slot := block.N.NewAlloca(types.Half)
	block.N.NewStore(constant.NewFloat(types.Half, float64(init)), slot)
	return &Float16{NativeType: types.Half, Value: slot, GoVal: init}
}

func (f *Float16) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, f.Value) }
func (f *Float16) Load(block *bc.BlockHolder) value.Value {
	return block.N.NewLoad(types.Half, f.Value)
}
func (f *Float16) Constant() constant.Constant {
	return constant.NewFloat(types.Half, float64(f.GoVal))
}
func (f *Float16) Slot() value.Value { return f.Value }
func (f *Float16) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.N.NewSIToFP(v, types.Half), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.N.NewFPTrunc(v, types.Half), nil
		case types.Half:
			return block.N.NewFPExt(v, types.Half), nil
		case types.Float:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
}
func (c *Float16) Type() types.Type         { return c.NativeType }
func (f *Float16) NativeTypeString() string { return "float16" }
