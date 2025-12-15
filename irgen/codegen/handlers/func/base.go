package funcs

import "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"

type FuncHandler struct {
	st *state.State
}

func NewFuncHandler(st *state.State) *FuncHandler {
	return &FuncHandler{
		st: st,
	}
}

var FuncHandlerInst *FuncHandler
