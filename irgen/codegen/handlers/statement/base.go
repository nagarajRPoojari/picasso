package statement

import "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"

// StatementHandler facilitates the translation of AST statement nodes into
// their corresponding IR representations.
type StatementHandler struct {
	st *state.State
}

// NewStatementHandler creates a constructor for StatementHandler, ensuring
// that it is properly bound to the compilation's lifecycle state.
func NewStatementHandler(st *state.State) *StatementHandler {
	return &StatementHandler{
		st: st,
	}
}

// StatementHandlerInst serves as a package-level singleton to provide
// global access to statement processing logic.
var StatementHandlerInst *StatementHandler
