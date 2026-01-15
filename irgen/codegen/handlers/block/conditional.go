package block

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

// processIfElseBlock orchestrates the generation of branching control flow in LLVM IR.
// It handles the evaluation of the boolean condition, creates discrete basic blocks
// for 'if' and 'else' execution paths, and ensures that all paths eventually
// synchronize at a common successor (phi-convergence) block.
//
// Key logic:
//   - Coerces the condition expression to a 1-bit integer (i1) via implicit casting.
//   - Automatically generates terminal jump instructions (Br) for blocks lacking a
//     explicit return or break, preventing malformed LLVM IR.
//   - Updates the provided BlockHolder to point to the 'end' block, allowing
//     subsequent statements to continue linearly.
//
// Note: endBlock accepted as paramas: it tells where to do correct phi convergence.
// This is needed in else-if ladder because endblocks are created by upper level blocks,
// so nested blocks need to know right endblock to converge
func (t *BlockHandler) processIfElseBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.IfStatement, endBlock *bc.BlockHolder) {
	ifBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	// Create end block only once (top-level)
	if endBlock == nil {
		endBlock = bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	}

	// condition
	res := t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(bh, st.Condition)
	cond := res.Load(bh)
	cond = t.st.TypeHandler.ImplicitIntCast(bh, cond, types.I1)

	var elseBlock *bc.BlockHolder

	if st.Alternate != nil {
		elseBlock = bc.NewBlockHolder(bh.V, fn.NewBlock(""))
		bh.N.NewCondBr(cond, ifBlock.N, elseBlock.N)
	} else {
		bh.N.NewCondBr(cond, ifBlock.N, endBlock.N)
	}

	// if block
	t.ProcessBlock(fn, ifBlock, st.Consequent.(ast.BlockStatement).Body)
	if ifBlock.N.Term == nil {
		ifBlock.N.NewBr(endBlock.N)
	}

	// else / else-if
	if st.Alternate != nil {
		switch alt := st.Alternate.(type) {
		case ast.BlockStatement:
			t.ProcessBlock(fn, elseBlock, alt.Body)
			if elseBlock.N.Term == nil {
				elseBlock.N.NewBr(endBlock.N)
			}

		// represents an else if ladder
		case ast.IfStatement:
			t.processIfElseBlock(fn, elseBlock, &alt, endBlock)
		}
	}

	bh.Update(endBlock.V, endBlock.N)
}
