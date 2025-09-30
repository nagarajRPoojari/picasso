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
//	entry – the current IR basic block before the if-statement
//	st    – the AST IfStatement to process
//
// Return:
//
//	A new IR basic block representing the merge point after the if-else.
func (t *BlockHandler) processIfElseBlock(fn *ir.Func, entry *ir.Block, st *ast.IfStatement) *ir.Block {
	ifBlock := fn.NewBlock("")
	elseBlock := fn.NewBlock("")
	endBlock := fn.NewBlock("")

	// condition
	res, safe := expression.ExpressionHandlerInst.ProcessExpression(entry, st.Condition)
	entry = safe

	casted, entry := t.st.TypeHandler.ImplicitTypeCast(entry, string(tf.BOOLEAN), res.Load(entry))
	cond := t.st.TypeHandler.BuildVar(entry, tf.Type(tf.BOOLEAN), casted)
	entry.NewCondBr(cond.Load(entry), ifBlock, elseBlock)

	// process consequent
	conseq := st.Consequent.(ast.BlockStatement)
	end := t.ProcessBlock(fn, ifBlock, conseq.Body)
	if end.Term == nil {
		end.NewBr(endBlock)
	}

	// process alternate
	alternate := st.Alternate.(ast.BlockStatement)
	end = t.ProcessBlock(fn, elseBlock, alternate.Body)
	if end.Term == nil {
		end.NewBr(endBlock)
	}

	return endBlock
}
