// Package class provides the logic for orchestrating User-Defined Type (UDT)
// lifecycle and object-oriented semantics within the LLVM IR generation phase.
//
// It handles the multi-pass process of class compilation:
//  1. Declaration: Registering opaque struct types and metadata containers.
//  2. Structural Definition: Calculating memory offsets, handling field inheritance,
//     and finalizing the LLVM struct layout.
//  3. Method Dispatch: Managing method signatures and resolving function pointers
//     for polymorphism and member access.
package class

import "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"

// ClassHandler manages the transformation of class declarations into
// LLVM-compatible structures and methods. It acts as the primary coordinator
// for the compiler's object model, provides methods to declare & define
// opaque & concrete classes
type ClassHandler struct {
	// st provides access to the global compiler state, including the
	// LLVM module, type registry, and identifier builder.
	st *state.State
}

// NewClassHandler creates a new handler instance with a shared reference
// to the compilation state. This link is essential for registering
// new types into the module-wide type system.
func NewClassHandler(state *state.State) *ClassHandler {
	return &ClassHandler{st: state}
}

// ClassHandlerInst is the global singleton instance used by the
// code generator to process class-related AST nodes during the
// various phases of the backend pipeline.
var ClassHandlerInst *ClassHandler
