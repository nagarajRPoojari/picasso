package block

import (
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/state"
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
