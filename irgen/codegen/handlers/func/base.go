// funcs package provides methods to declare & define function blocks.
package funcs

import (
	"github.com/nagarajRPoojari/niyama/irgen/codegen/contract"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
)

// FuncHandler encapsulates the state required to generate IR
// for function declaration & definition
type FuncHandler struct {
	st *state.State
	m  contract.Mediator
}

// NewFuncHandler initializes the function handler
func NewFuncHandler(st *state.State, m contract.Mediator) *FuncHandler {
	return &FuncHandler{
		st: st,
		m:  m,
	}
}
