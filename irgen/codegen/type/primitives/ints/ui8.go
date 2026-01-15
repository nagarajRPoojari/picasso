package ints

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
)

// Int8 hold 1 byte of information
type UInt8 struct {
	NativeType *types.IntType
	Value      value.Value
}

func NewUInt8Var(block *bc.BlockHolder, init int8) *UInt8 {
	slot := block.N.NewAlloca(types.I8)
	block.N.NewStore(constant.NewInt(types.I8, int64(init)), slot)
	return &UInt8{NativeType: types.I8, Value: slot}
}

func (i *UInt8) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, i.Value) }
func (i *UInt8) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I8, i.Value) }
func (i *UInt8) Constant() constant.Constant                 { return constant.NewInt(types.I8, int64(0)) }
func (i *UInt8) Slot() value.Value                           { return i.Value }
func (c *UInt8) Type() types.Type                            { return c.NativeType }
func (i *UInt8) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		// truncate to 1 byte if byte size if greater
		if t.BitSize > 8 {
			return block.N.NewTrunc(v, types.I8), nil
		} else if t.BitSize < 8 {
			return block.N.NewSExt(v, types.I8), nil
		}
		return v, nil
	case *types.FloatType:
		return block.N.NewFPToSI(v, types.I8), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int7", v))
	}
}
func (f *UInt8) NativeTypeString() string { return "uint8" }
