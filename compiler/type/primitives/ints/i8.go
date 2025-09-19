package ints

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Int8 hold 1 byte of information
type Int8 struct {
	NativeType *types.IntType
	Value      value.Value
	GoVal      int8
}

func NewInt8Var(block *ir.Block, init int8) *Int8 {
	slot := block.NewAlloca(types.I8)
	block.NewStore(constant.NewInt(types.I8, int64(init)), slot)
	return &Int8{NativeType: types.I8, Value: slot, GoVal: init}
}

func (i *Int8) Update(block *ir.Block, v value.Value) { block.NewStore(v, i.Value) }
func (i *Int8) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I8, i.Value) }
func (i *Int8) Constant() constant.Constant           { return constant.NewInt(types.I8, int64(i.GoVal)) }
func (i *Int8) Slot() value.Value                     { return i.Value }
func (c *Int8) Type() types.Type                      { return c.NativeType }
func (i *Int8) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		// truncate to 1 byte if byte size if greater
		if t.BitSize > 8 {
			return block.NewTrunc(v, types.I8), nil
		} else if t.BitSize < 8 {
			return block.NewSExt(v, types.I8), nil
		}
		return v, nil
	case *types.FloatType:
		return block.NewFPToSI(v, types.I8), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int7", v))
	}
}
func (f *Int8) NativeTypeString() string { return "int8" }
