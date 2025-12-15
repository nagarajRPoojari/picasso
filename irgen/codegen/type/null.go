package typedef

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// NullVar represents a null variable (used for uninitialized or placeholder values).
type NullVar struct {
	typ *types.PointerType
}

func NewNullVar(typ *types.PointerType) *NullVar {
	return &NullVar{typ: typ}
}

func (n *NullVar) Update(block *bc.BlockHolder, v value.Value) {
	// Null cannot be updated, silently ignore or panic depending on your design choice
	// panic("cannot update null variable")
}

func (n *NullVar) Load(block *bc.BlockHolder) value.Value {
	// Always return a null constant for its type
	return constant.NewNull(n.typ)
}

func (n *NullVar) Constant() constant.Constant {
	return constant.NewNull(n.typ)
}

func (n *NullVar) Slot() value.Value {
	// Null has no slot, return nil
	return nil
}

func (n *NullVar) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	// If v already matches the type, return as is
	if v.Type().Equal(n.typ) {
		return v, nil
	}
	// Else try casting null (only valid for pointer types usually)
	return constant.NewNull(n.typ), nil
}

func (n *NullVar) Type() types.Type {
	return n.typ
}
func (f *NullVar) NativeTypeString() string { return "null" }
