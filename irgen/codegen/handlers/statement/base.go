package statement

import (
	"github.com/nagarajRPoojari/picasso/irgen/codegen/contract"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
)

// StatementHandler facilitates the translation of AST statement nodes into
// their corresponding IR representations.
type StatementHandler struct {
	st *state.State
	m  contract.Mediator
}

// NewStatementHandler creates a constructor for StatementHandler, ensuring
// that it is properly bound to the compilation's lifecycle state.
func NewStatementHandler(st *state.State, m contract.Mediator) *StatementHandler {
	return &StatementHandler{
		st: st,
		m:  m,
	}
}
