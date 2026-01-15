package ints

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
)

// Int16 store 2 bytes of information
type Int16 struct {
	NativeType *types.IntType
	Value      value.Value
	GoVal      int16
}

func NewInt16Var(block *bc.BlockHolder, init int16) *Int16 {
	slot := block.N.NewAlloca(types.I16)
	block.N.NewStore(constant.NewInt(types.I16, int64(init)), slot)
	return &Int16{NativeType: types.I16, Value: slot, GoVal: init}
}

func (i *Int16) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, i.Value) }
func (i *Int16) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I16, i.Value) }
func (i *Int16) Constant() constant.Constant                 { return constant.NewInt(types.I16, int64(i.GoVal)) }
func (i *Int16) Slot() value.Value                           { return i.Value }
func (c *Int16) Type() types.Type                            { return c.NativeType }
func (i *Int16) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
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
func (f *Int16) NativeTypeString() string { return "int16" }
