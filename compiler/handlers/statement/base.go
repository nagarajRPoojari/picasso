package statement

import "github.com/nagarajRPoojari/x-lang/compiler/handlers/state"

type StatementHandler struct {
	st *state.State
}

func NewStatementHandler(st *state.State) *StatementHandler {
	return &StatementHandler{
		st: st,
	}
}

var StatementHandlerInst *StatementHandler
