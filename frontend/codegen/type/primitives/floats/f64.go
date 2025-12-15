package floats

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/niyama/frontend/error"
)

// Float64 stores 8 byte floating point value
type Float64 struct {
	NativeType *types.FloatType
	Value      value.Value
	GoVal      float64
}

func NewFloat64Var(block *bc.BlockHolder, init float64) *Float64 {
	slot := block.N.NewAlloca(types.Double)
	block.N.NewStore(constant.NewFloat(types.Double, init), slot)
	return &Float64{NativeType: types.Double, Value: slot, GoVal: init}
}

func (f *Float64) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, f.Value) }
func (f *Float64) Load(block *bc.BlockHolder) value.Value {
	return block.N.NewLoad(types.Double, f.Value)
}
func (f *Float64) Constant() constant.Constant { return constant.NewFloat(types.Double, f.GoVal) }
func (f *Float64) Slot() value.Value           { return f.Value }
func (c *Float64) Type() types.Type            { return c.NativeType }
func (f *Float64) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.N.NewSIToFP(v, types.Double), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Float:
			return block.N.NewFPExt(v, types.Double), nil
		case types.Half:
			return block.N.NewFPExt(v, types.Double), nil
		case types.Double:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to float64", v))
}
func (f *Float64) NativeTypeString() string { return "float64" }
