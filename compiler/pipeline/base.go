package pipeline

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/c"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
)

type Pipeline struct {
	st   *state.State
	tree ast.BlockStatement
}

func NewPipeline(st *state.State, tree ast.BlockStatement) *Pipeline {
	return &Pipeline{st: st, tree: tree}
}

func (t *Pipeline) Run() {
	t.ImportModules()
	t.PredeclareClasses()
	t.DeclareVars()
	t.DeclareFuncs()
	t.DefineClasses()
	t.DefineMain()

	insertYields(t.st.Module)
}

func insertYields(m *ir.Module) {
	// Define or get the yield function
	yieldFunc := c.NewInterface(m).Funcs[c.SELF_YIELD]

	for _, fn := range m.Funcs {
		// Skip the yield function itself
		if fn.Name() == yieldFunc.Name() {
			continue
		}

		for _, blk := range fn.Blocks {
			var newInsts []ir.Instruction

			for _, inst := range blk.Insts {
				newInsts = append(newInsts, inst)

				// === Heuristic: insert yield after every call ===
				switch inst.(type) {
				case *ir.InstCall:
					// Insert a call to yield_func
					newInsts = append(newInsts, ir.NewCall(yieldFunc))
				}
			}

			// Replace block instructions with new sequence
			blk.Insts = newInsts
		}
	}
}
