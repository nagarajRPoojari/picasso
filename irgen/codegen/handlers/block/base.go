/*
Package block provides the logic for handling scoped code execution.
It manages the entry and exit of lexical blocks, ensuring that local
variables and control-flow instructions are correctly bound to the
current execution context.
*/
package block

import (
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
)

// BlockHandler is responsible for orchestrating the generation of IR
// within a specific lexical scope. It interacts with the global
// compiler state to manage symbol visibility and stack allocation.
type BlockHandler struct {
	// st represents the shared compilation state, including
	// the current LLVM builder and symbol tables.
	st *state.State
}

// NewBlockHandler initializes a handler with a reference to the
// compiler's current state. This allows the handler to modify
// the IR module and track scope-specific metadata.
func NewBlockHandler(st *state.State) *BlockHandler {
	return &BlockHandler{
		st: st,
	}
}

// BlockHandlerInst is a global singleton instance of the BlockHandler.
// It is typically initialized during the setup of the code generator
// to provide a consistent entry point for block-level IR emission.
var BlockHandlerInst *BlockHandler
