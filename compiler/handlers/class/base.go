package class

import "github.com/nagarajRPoojari/x-lang/compiler/handlers/state"

type ClassHandler struct {
	st *state.State
}

func NewClassHandler(state *state.State) *ClassHandler {
	return &ClassHandler{st: state}
}

var ClassHandlerInst *ClassHandler
