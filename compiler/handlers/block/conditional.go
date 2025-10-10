package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
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

	elseBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	// condition
	res := expression.ExpressionHandlerInst.ProcessExpression(bh, st.Condition)

	casted := t.st.TypeHandler.ImplicitTypeCast(bh, string(tf.BOOLEAN), res.Load(bh))

	cond := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.BOOLEAN), casted)
	bh.N.NewCondBr(cond.Load(bh), ifBlock.N, elseBlock.N)

	// process consequent
	conseq := st.Consequent.(ast.BlockStatement)
	t.ProcessBlock(fn, ifBlock, conseq.Body)
	if ifBlock.N.Term == nil {
		ifBlock.N.NewBr(endBlock.N)
	}

	// process alternate
	alternate := st.Alternate.(ast.BlockStatement)
	t.ProcessBlock(fn, elseBlock, alternate.Body)
	if elseBlock.N.Term == nil {
		elseBlock.N.NewBr(endBlock.N)
	}

	bh.Update(endBlock.V, endBlock.N)
}
