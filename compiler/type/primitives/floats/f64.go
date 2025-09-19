package floats

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Float64 stores 8 byte floating point value
type Float64 struct {
	NativeType *types.FloatType
	Value      value.Value
	GoVal      float64
}

func NewFloat64Var(block *ir.Block, init float64) *Float64 {
	slot := block.NewAlloca(types.Double)
	block.NewStore(constant.NewFloat(types.Double, init), slot)
	return &Float64{NativeType: types.Double, Value: slot, GoVal: init}
}

func (f *Float64) Update(block *ir.Block, v value.Value) { block.NewStore(v, f.Value) }
func (f *Float64) Load(block *ir.Block) value.Value      { return block.NewLoad(types.Double, f.Value) }
func (f *Float64) Constant() constant.Constant           { return constant.NewFloat(types.Double, f.GoVal) }
func (f *Float64) Slot() value.Value                     { return f.Value }
func (c *Float64) Type() types.Type                      { return c.NativeType }
func (f *Float64) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Double), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Float:
			return block.NewFPExt(v, types.Double), nil
		case types.Half:
			return block.NewFPExt(v, types.Double), nil
		case types.Double:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to float64", v))
}
func (f *Float64) NativeTypeString() string { return "float64" }
