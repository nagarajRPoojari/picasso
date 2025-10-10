package ints

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

type Int32 struct {
	NativeType *types.IntType
	Value      value.Value
	GoVal      int32
}

func NewInt32Var(block *bc.BlockHolder, init int32) *Int32 {
	slot := block.N.NewAlloca(types.I32)
	block.N.NewStore(constant.NewInt(types.I32, int64(init)), slot)
	return &Int32{NativeType: types.I32, Value: slot, GoVal: init}
}

func (i *Int32) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, i.Value) }
func (i *Int32) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I32, i.Value) }
func (i *Int32) Constant() constant.Constant                 { return constant.NewInt(types.I32, int64(i.GoVal)) }
func (i *Int32) Slot() value.Value                           { return i.Value }
func (i *Int32) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 32 {
			return block.N.NewTrunc(v, types.I32), nil
		} else if t.BitSize < 32 {
			return block.N.NewSExt(v, types.I32), nil
		}
		return v, nil
	case *types.FloatType:
		return block.N.NewFPToSI(v, types.I32), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int32", v))
	}
}
func (c *Int32) Type() types.Type         { return c.NativeType }
func (f *Int32) NativeTypeString() string { return "int32" }
