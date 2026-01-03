// Package interfaceh provides the logic for orchestrating Interface lifecycle
// and contract-based polymorphism within the LLVM IR generation phase.
//
// It manages a three-pass compilation strategy for interfaces:
//  1. Declaration: Creating named opaque structs to support forward references.
//  2. Method Registration: Generating function prototypes and calculating
//     signature hashes for implementation validation.
package interfaceh

import (
	"github.com/nagarajRPoojari/niyama/irgen/codegen/contract"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
)

// InterfaceHandler manages the transformation of Niyama interface declarations
// into LLVM-compatible virtual tables. It ensures that method signatures
// are consistently hashed and indexed, allowing implementing classes to
// satisfy the interface's structural requirements.
type InterfaceHandler struct {
	// st provides access to the global compiler state, including the
	// LLVM module, type registry, and cross-package symbol tables.
	st *state.State

	m contract.Mediator
}

// NewInterfaceHandler initializes a handler with the shared compilation state.
// This allows the handler to register symbolic function definitions and
// concrete types into the module's global scope.
func NewInterfaceHandler(state *state.State, m contract.Mediator) *InterfaceHandler {
	return &InterfaceHandler{st: state, m: m}
}
