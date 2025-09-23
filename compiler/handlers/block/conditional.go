package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *BlockHandler) processIfElseBlock(fn *ir.Func, entry *ir.Block, st *ast.IfStatement) *ir.Block {
	ifBlock := fn.NewBlock("")
	elseBlock := fn.NewBlock("")
	endBlock := fn.NewBlock("")

	// condition
	res := expression.ExpressionHandlerInst.ProcessExpression(entry, st.Condition)
	casted := t.st.TypeHandler.CastToType(entry, string(tf.BOOLEAN), res.Load(entry))
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
