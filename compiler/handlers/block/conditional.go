package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
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
func (t *BlockHandler) processIfElseBlock(fn *ir.Func, bh tf.BlockHolder, st *ast.IfStatement) tf.BlockHolder {
	ifBlock := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}

	elseBlock := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}
	endBlock := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}

	// condition
	res, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, st.Condition)
	bh = safe

	casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, string(tf.BOOLEAN), res.Load(bh.N))
	bh.N = safeN

	cond := t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.BOOLEAN), casted)
	bh.N.NewCondBr(cond.Load(bh.N), ifBlock.N, elseBlock.N)

	// process consequent
	conseq := st.Consequent.(ast.BlockStatement)
	end := t.ProcessBlock(fn, ifBlock, conseq.Body)
	if end.N.Term == nil {
		end.N.NewBr(endBlock.N)
	}

	// process alternate
	alternate := st.Alternate.(ast.BlockStatement)
	end = t.ProcessBlock(fn, elseBlock, alternate.Body)
	if end.N.Term == nil {
		end.N.NewBr(endBlock.N)
	}

	return endBlock
}
