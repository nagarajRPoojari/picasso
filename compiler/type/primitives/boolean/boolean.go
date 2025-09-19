package boolean

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/compiler/utils"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Boolean hold single bit of information
type Boolean struct {
	NativeType *types.IntType
	Value      value.Value // pointer to i1
	GoVal      bool
}

func NewBooleanVar(block *ir.Block, init bool) *Boolean {
	slot := block.NewAlloca(types.I1)
	block.NewStore(constant.NewInt(types.I1, utils.BtoI(init)), slot)
	return &Boolean{NativeType: types.I1, Value: slot, GoVal: init}
}
func (b *Boolean) Update(block *ir.Block, v value.Value) { block.NewStore(v, b.Value) }
func (b *Boolean) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I1, b.Value) }
func (b *Boolean) Constant() constant.Constant           { return constant.NewInt(types.I1, utils.BtoI(b.GoVal)) }
func (b *Boolean) Slot() value.Value                     { return b.Value }
func (c *Boolean) Type() types.Type                      { return c.NativeType }
func (b *Boolean) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		// cast any int to i1
		if v.Type().(*types.IntType).BitSize == 1 {
			return v, nil
		}
		// if incoming int type is other than types.I1, truncate it to single bit
		return block.NewTrunc(v, types.I1), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to booolean", v))
	}
}
func (f *Boolean) NativeTypeString() string { return "boolean" }
