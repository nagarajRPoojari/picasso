package ints

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
)

// UInt16 store 2 bytes of information
type UInt16 struct {
	NativeType *types.IntType
	Value      value.Value
	GoVal      int16
}

func NewUInt16Var(block *bc.BlockHolder, init int16) *UInt16 {
	slot := block.N.NewAlloca(types.I16)
	block.N.NewStore(constant.NewInt(types.I16, int64(init)), slot)
	return &UInt16{NativeType: types.I16, Value: slot, GoVal: init}
}

func (i *UInt16) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, i.Value) }
func (i *UInt16) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I16, i.Value) }
func (i *UInt16) Constant() constant.Constant                 { return constant.NewInt(types.I16, int64(i.GoVal)) }
func (i *UInt16) Slot() value.Value                           { return i.Value }
func (c *UInt16) Type() types.Type                            { return c.NativeType }
func (i *UInt16) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 16 {
			return block.N.NewTrunc(v, types.I16), nil
		} else if t.BitSize < 16 {
			return block.N.NewSExt(v, types.I16), nil
		}
		return v, nil
	case *types.FloatType:
		return block.N.NewFPToSI(v, types.I16), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int16", v))

	}
}
func (f *UInt16) NativeTypeString() string { return "uint16" }
