package block

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

// processIfElseBlock generates LLVM IR for an if-else statement.
// It builds condition evaluation, branching, and merges control flow
// into a common end block.
//
// Params:
//
//	fn    – the LLVM function being built
//	bh – the current IR basic block before the if-statement
//	st    – the AST IfStatement to process
//
// Return:
//
//	A new IR basic block representing the merge point after the if-else.
func (t *BlockHandler) processIfElseBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.IfStatement) {
	ifBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	// condition
	res := expression.ExpressionHandlerInst.ProcessExpression(bh, st.Condition)
	cond := res.Load(bh)
	cond = t.st.TypeHandler.ImplicitIntCast(bh, cond, types.I1)

	// handle else branch
	var elseBlock *bc.BlockHolder
	if st.Alternate != nil {
		elseBlock = bc.NewBlockHolder(bh.V, fn.NewBlock(""))
		bh.N.NewCondBr(cond, ifBlock.N, elseBlock.N)
	} else {
		bh.N.NewCondBr(cond, ifBlock.N, endBlock.N)
	}

	// process if block
	if st.Consequent != nil {
		if conseq, ok := st.Consequent.(ast.BlockStatement); ok {
			t.ProcessBlock(fn, ifBlock, conseq.Body)
		}
	}
	if ifBlock.N.Term == nil {
		ifBlock.N.NewBr(endBlock.N)
	}

	// process else block if exists
	if elseBlock != nil && st.Alternate != nil {
		if alt, ok := st.Alternate.(ast.BlockStatement); ok {
			t.ProcessBlock(fn, elseBlock, alt.Body)
		}
		if elseBlock.N.Term == nil {
			elseBlock.N.NewBr(endBlock.N)
		}
	}

	bh.Update(endBlock.V, endBlock.N)
}
