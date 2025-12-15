package bc

import "github.com/llir/llvm/ir"

type VarBlock struct {
	*ir.Block
}

type BlockHolder struct {
	V VarBlock
	N *ir.Block
}

func NewBlockHolder(v VarBlock, n *ir.Block) *BlockHolder {
	return &BlockHolder{
		V: v,
		N: n,
	}
}

func (t *BlockHolder) Update(v VarBlock, n *ir.Block) {
	t.N = n
	t.V = v
}
