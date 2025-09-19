package pipeline

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
)

type Pipeline struct {
	tree *ast.BlockStatement
	st   *state.State
}

func NewPipeline(st *state.State, tree *ast.BlockStatement) *Pipeline {
	return &Pipeline{st: st, tree: tree}
}

func (t *Pipeline) Run() {
	t.ImportModules()
	t.PredeclareClasses()
	t.DeclareVars()
	t.DeclareFuncs()
	t.DefineClasses()
	t.DefineMain()
}
