package block

import (
	"github.com/nagarajRPoojari/x-lang/generator/handlers/state"
)

type BlockHandler struct {
	st *state.State
}

func NewBlockHandler(st *state.State) *BlockHandler {
	return &BlockHandler{
		st: st,
	}
}

var BlockHandlerInst *BlockHandler
