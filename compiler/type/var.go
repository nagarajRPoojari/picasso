package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Var is a holder for all native variables
// It holds mutable slot suporting runtime update and load operations
type Var interface {
	Update(block *ir.Block, v value.Value)
	Load(block *ir.Block) value.Value

	// mutable slot holding the value
	Slot() value.Value

	// Cast casts given value to self type if possible
	Cast(block *ir.Block, v value.Value) (value.Value, error)

	// Type returns llvm compatibe type
	Type() types.Type
}
