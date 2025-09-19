package funcs

import "github.com/nagarajRPoojari/x-lang/compiler/handlers/state"

type FuncHandler struct {
	st *state.State
}

func NewFuncHandler(st *state.State) *FuncHandler {
	return &FuncHandler{
		st: st,
	}
}

var FuncHandlerInst *FuncHandler
