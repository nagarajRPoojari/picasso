package pipeline

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/state"
)

type Pipeline struct {
	st   *state.State
	tree ast.BlockStatement
}

func NewPipeline(st *state.State, tree ast.BlockStatement) *Pipeline {
	return &Pipeline{st: st, tree: tree}
}

func (t *Pipeline) Run() {
	t.Register()
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
				// If this is a call instruction, insert yield BEFORE it
				switch inst.(type) {
				case *ir.InstCall:
					newInsts = append(newInsts, ir.NewCall(yieldFunc))
				}

				// Then append the original instruction
				newInsts = append(newInsts, inst)
			}

			blk.Insts = newInsts
		}
	}
}
