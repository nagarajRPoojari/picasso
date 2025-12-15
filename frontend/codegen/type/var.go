package typedef

import (
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

// Var represents a container for all native variables in the system.
// It provides a mutable slot that supports runtime update and load operations.
type Var interface {
	// Update stores the given value into the variable's slot.
	Update(block *bc.BlockHolder, v value.Value)

	// Load retrieves the current value from the variable's slot.
	Load(block *bc.BlockHolder) value.Value

	// Slot returns the underlying mutable storage for this variable.
	// Must be pointer to value type.
	Slot() value.Value

	// Cast attempts to convert the given value into this variable’s type.
	// @todo: deprecate
	Cast(block *bc.BlockHolder, v value.Value) (value.Value, error)

	// Type returns the LLVM-compatible type of this variable.
	Type() types.Type

	// NativeTypeString returns the native type name of this variable.
	// Example: "Math", "string", "int".
	NativeTypeString() string
}
