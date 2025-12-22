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

// Declare emits llvm declaration instructions.
// sourcePkg indciates where those type/functions definitions comes from.
// it is needed to determine its fully qualified name.
// e.g, using "os/io"; declaration should be os.io.ABC etc..
func (t *Pipeline) Declare(sourcePkg state.PackageEntry) {
	t.predeclareInterfraces(sourcePkg)

	t.predeclareClasses(sourcePkg)

	t.declareInterfaceFields(sourcePkg)
	t.declareInterfaceFuncs(sourcePkg)

	t.declareClassFields(sourcePkg)
	t.declareClassFuncs(sourcePkg)
}

// Definitions are called for own module which emits definition
// instructions in llvm.
func (t *Pipeline) Define() {
	t.defineClasses()
	t.defineMain()
}

func (t *Pipeline) Optimize() {
	// inserts safepoints based on some predefined heuristics.
	t.insertYields()
}

// Run will be called only for own module which does both
// declaration & definition. It is expected that all imported
// types/funcs are already declared.
func (t *Pipeline) Run(sourcePkg state.PackageEntry) {
	t.Register()
	t.Declare(sourcePkg)
	t.Define()
	t.Optimize()
}
