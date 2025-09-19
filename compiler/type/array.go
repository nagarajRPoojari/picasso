package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

type Array struct {
	NativeType *types.ArrayType
	EleType    *types.Type
	Value      value.Value
	Goval      int
}

func NewArray(block *ir.Block, t types.Type, size int) *Array {
	arrType := types.NewArray(uint64(size), t)
	slot := block.NewAlloca(arrType)

	return &Array{
		NativeType: arrType,
		Value:      slot,
		EleType:    &t,
	}
}

func (b *Array) Update(block *ir.Block, v value.Value) { block.NewStore(v, b.Value) }
func (b *Array) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I1, b.Value) }
func (b *Array) Slot() value.Value                     { return b.Value }
func (c *Array) Type() types.Type                      { return c.NativeType }
func (b *Array) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to array", v))
}

func (a *Array) Index(block *ir.Block, indices ...value.Value) value.Value {
	zero := constant.NewInt(types.I32, 0)
	idxs := append([]value.Value{zero}, indices...)
	return block.NewGetElementPtr(a.NativeType, a.Value, idxs...)
}

func (a *Array) StoreAt(block *ir.Block, v value.Value, indices ...value.Value) {
	ptr := a.Index(block, indices...)
	block.NewStore(v, ptr)
}

func (a *Array) LoadAt(block *ir.Block, indices ...value.Value) value.Value {
	ptr := a.Index(block, indices...)
	return block.NewLoad((*a.EleType), ptr)
}
