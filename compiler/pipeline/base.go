package pipeline

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
)

type Pipeline struct {
	packages map[string]ast.BlockStatement
	st       *state.State

	tree ast.BlockStatement
}

func NewPipeline(st *state.State, packages map[string]ast.BlockStatement) *Pipeline {
	return &Pipeline{st: st, packages: packages}
}

func (t *Pipeline) Run() {
	t.tree = t.BuildUniModule()
	t.ImportModules()
	t.PredeclareClasses()
	t.DeclareVars()
	t.DeclareFuncs()
	t.DefineClasses()
	t.DefineMain()
}
