package floats

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Float32 stores 4 byte floating point value
type Float32 struct {
	NativeType *types.FloatType
	Value      value.Value
	GoVal      float32
}

func NewFloat32Var(block *bc.BlockHolder, init float32) *Float32 {
	slot := block.N.NewAlloca(types.Float)
	block.N.NewStore(constant.NewFloat(types.Float, float64(init)), slot)
	return &Float32{NativeType: types.Float, Value: slot, GoVal: init}
}

func (f *Float32) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, f.Value) }
func (f *Float32) Load(block *bc.BlockHolder) value.Value {
	return block.N.NewLoad(types.Float, f.Value)
}
func (f *Float32) Constant() constant.Constant {
	return constant.NewFloat(types.Float, float64(f.GoVal))
}
func (f *Float32) Slot() value.Value { return f.Value }
func (c *Float32) Type() types.Type  { return c.NativeType }
func (f *Float32) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.N.NewSIToFP(v, types.Float), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.N.NewFPTrunc(v, types.Float), nil
		case types.Half:
			return block.N.NewFPExt(v, types.Float), nil
		case types.Float:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to float32", v))
}
func (f *Float32) NativeTypeString() string { return "float32" }
