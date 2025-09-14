package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type Null struct {
	Target types.Type
	Value  *ir.InstAlloca
}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Load(block *ir.Block) value.Value {
	return block.NewLoad(n.Target, n.Value)
}

func (n *Null) Update(block *ir.Block, v value.Value) {
	if !v.Type().Equal(n.Target) {
		panic(fmt.Sprintf("cannot assign %s to null of type %s", v.Type(), n.Target))
	}
	block.NewStore(v, n.Value)
}

func (n *Null) Type() types.Type {
	return n.Target
}

func (n *Null) Cast(block *ir.Block, v value.Value) value.Value {
	panic(fmt.Sprintf("cannot cast %s to null type %s", v.Type(), n.Target))
}

func (n *Null) GoValue() interface{} {
	return nil
}
