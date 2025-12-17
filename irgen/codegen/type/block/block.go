package bc

import "github.com/llir/llvm/ir"

// wrapping with VarBlock to strictly separate var declaration
// instructions & others.
type VarBlock struct {
	*ir.Block
}

// BlockHolder holds *ir.Block instance.
type BlockHolder struct {
	// should only be used for var declaration
	V VarBlock
	// should be used for rest of the instructions
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
