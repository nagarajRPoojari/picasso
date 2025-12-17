package pipeline

import (
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
)

type Pipeline struct {
	st   *state.State
	tree ast.BlockStatement
}

func NewPipeline(st *state.State, tree ast.BlockStatement) *Pipeline {
	return &Pipeline{st: st, tree: tree}
}

func (t *Pipeline) Register() {
	t.registerTypes()
}

func (t *Pipeline) Declare() {
	t.predeclareClasses()
	t.declareVars()
	t.declareFuncs()
}

func (t *Pipeline) Define() {
	t.defineClasses()
	t.defineMain()
}

func (t *Pipeline) Optimize() {
	t.insertYields()
}

func (t *Pipeline) Run() {
	t.Register()
	t.Declare()
	t.Define()
	t.Optimize()
}
