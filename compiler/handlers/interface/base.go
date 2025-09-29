package _interface

import "github.com/nagarajRPoojari/x-lang/compiler/handlers/state"

type InterfaceHandler struct {
	st *state.State
}

func NewClassHandler(state *state.State) *InterfaceHandler {
	return &InterfaceHandler{st: state}
}

var InterfaceHandlerInst *InterfaceHandler
