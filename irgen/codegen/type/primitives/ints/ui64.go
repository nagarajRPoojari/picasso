package ints

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/niyama/irgen/error"
)

// UInt64 stores 8 bytes int value
type UInt64 struct {
	NativeType *types.IntType
	Value      value.Value
	GoVal      int64
}

func NewUInt64Var(block *bc.BlockHolder, init int64) *UInt64 {
	slot := block.N.NewAlloca(types.I64)
	block.N.NewStore(constant.NewInt(types.I64, init), slot)
	return &UInt64{NativeType: types.I64, Value: slot, GoVal: init}
}

func (i *UInt64) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, i.Value) }
func (i *UInt64) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I64, i.Value) }
func (i *UInt64) Constant() constant.Constant                 { return constant.NewInt(types.I64, i.GoVal) }
func (i *UInt64) Slot() value.Value                           { return i.Value }
func (i *UInt64) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 64 {
			return block.N.NewTrunc(v, types.I64), nil
		} else if t.BitSize < 64 {
			return block.N.NewSExt(v, types.I64), nil
		}
		return v, nil
	case *types.FloatType:
		return block.N.NewFPToSI(v, types.I64), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
	}
}
func (c *UInt64) Type() types.Type         { return c.NativeType }
func (f *UInt64) NativeTypeString() string { return "uint64" }
