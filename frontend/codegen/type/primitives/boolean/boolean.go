package boolean

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/utils"
	errorsx "github.com/nagarajRPoojari/niyama/frontend/error"
)

// Boolean hold single bit of information
type Boolean struct {
	NativeType *types.IntType
	Value      value.Value // pointer to i1
	GoVal      bool
}

func NewBooleanVar(bh *bc.BlockHolder, init bool) *Boolean {
	block := bh.N
	slot := block.NewAlloca(types.I1)
	block.NewStore(constant.NewInt(types.I1, utils.BtoI(init)), slot)
	return &Boolean{NativeType: types.I1, Value: slot, GoVal: init}
}
func (b *Boolean) Update(block *bc.BlockHolder, v value.Value) { block.N.NewStore(v, b.Value) }
func (b *Boolean) Load(block *bc.BlockHolder) value.Value      { return block.N.NewLoad(types.I1, b.Value) }
func (b *Boolean) Constant() constant.Constant                 { return constant.NewInt(types.I1, utils.BtoI(b.GoVal)) }
func (b *Boolean) Slot() value.Value                           { return b.Value }
func (c *Boolean) Type() types.Type                            { return c.NativeType }
func (b *Boolean) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize == 1 {
			return v, nil // already i1
		}
		zero := constant.NewInt(t, 0)
		return block.N.NewICmp(enum.IPredNE, v, zero), nil

	case *types.FloatType:
		zero := constant.NewFloat(t, 0.0)
		return block.N.NewFCmp(enum.FPredONE, v, zero), nil

	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
	}
}
func (f *Boolean) NativeTypeString() string { return "boolean" }
