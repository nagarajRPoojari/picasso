package block

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// processForBlock implements the IR lowering for 'foreach' style range loops.
// It transforms a high-level range expression into a classic four-block
// loop structure: Header (Condition), Body, Increment (Latch), and Exit.
//
// Key logic:
//   - Range Lowering: Evaluates the 'Lower' and 'Upper' expressions and casts
//     them to standard integer types for comparison logic.
//   - Iterator Management: Allocates a stack slot for the loop variable and
//     registers it in a new lexical scope so it is accessible within the body.
//   - Control Flow: Manages the 'Loopend' stack to support 'break' statements
//     inside the loop body, ensuring they jump to the correct exit block.
//   - Increment Logic: Automatically generates the i++ logic and branches
//     back to the header to re-evaluate the loop invariant.
func (t *BlockHandler) processForBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.ForeachStatement) {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	lowerExpr := st.Iterable.(ast.RangeExpression).Lower
	upperExpr := st.Iterable.(ast.RangeExpression).Upper

	// start initializing index variable from lower bound expression.
	// core logic is to put initialization in current block & create
	// a new block for condition & incrementing, so that i can branch back to it after
	// executing loop body.
	// Note: index variables are kept i64, need to @verify this.
	indexVal := t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(bh, lowerExpr)
	casted := t.st.TypeHandler.ImplicitTypeCast(bh, tf.INT, indexVal.Load(bh))
	indexVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	iPtr := indexVal.Slot()
	t.st.Vars.AddNewVar(st.Value, indexVal)

	upperVal := t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(bh, upperExpr)
	casted = t.st.TypeHandler.ImplicitTypeCast(bh, tf.INT, upperVal.Load(bh))
	upperVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	loopCond := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopBody := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopInc := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopEnd := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	bh.N.NewBr(loopCond.N)

	// condition need to be done in new block so that i can iterate back to it.
	// Note that it follows [a, b) logic, (i.e, excluding right bound)
	iVal := loopCond.N.NewLoad(types.I64, iPtr)
	cond := loopCond.N.NewICmp(enum.IPredSLT, iVal, upperVal.Load(loopCond))
	loopCond.N.NewCondBr(cond, loopBody.N, loopEnd.N)

	// loop blocks need to be appended to a temporary stack to remove 'break' statements
	// with respect to last pushed loop block
	t.st.Loopend = append(t.st.Loopend, state.LoopEntry{End: loopEnd})
	t.ProcessBlock(fn, loopBody, st.Body)
	t.st.Loopend = t.st.Loopend[:len(t.st.Loopend)-1]

	if loopBody.N.Term == nil {
		loopBody.N.NewBr(loopInc.N)
	}

	iVal2 := loopInc.N.NewLoad(types.I64, iPtr)
	iNext := loopInc.N.NewAdd(iVal2, constant.NewInt(types.I64, 1))
	loopInc.N.NewStore(iNext, iPtr)
	loopInc.N.NewBr(loopCond.N)

	bh.Update(loopEnd.V, loopEnd.N)
}

// processWhileBlock generates the LLVM IR representation for a while-loop construct.
// It establishes a cyclic control flow graph by partitioning the loop into
// three distinct basic blocks: a condition header, the loop body, and a
// post-loop exit block.
//
// Key logic:
//   - Header Branching: Creates a dedicated 'condBlock' to re-evaluate the
//     boolean condition at the start of every iteration.
//   - Scope Integrity: Manages a fresh variable block to isolate local
//     declarations defined within the loop body.
//   - Break Support: Pushes the 'endBlock' onto the Loopend stack, enabling
//     nested statements to resolve the correct jump target for 'break' commands.
//   - Back-edge Generation: Automatically injects an unconditional branch from
//     the end of the body back to the condition header, ensuring the loop persists.
func (t *BlockHandler) processWhileBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.WhileStatement) {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	condEntry := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	bodyBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	bh.N.NewBr(condEntry.N)

	// I'll modify the condition entry later, to loop back to condition entry i need to
	// store its state.
	copyOfCondEntry := bc.NewBlockHolder(condEntry.V, condEntry.N)

	res := t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(condEntry, st.Condition)

	condBlock := condEntry
	cond := res.Load(condBlock)
	cond = t.st.TypeHandler.ImplicitIntCast(condBlock, cond, types.I1)

	condBlock.N.NewCondBr(cond, bodyBlock.N, endBlock.N)

	bh.Update(bodyBlock.V, bodyBlock.N)

	// loop blocks need to be appended to a temporary stack to remove 'break' statements
	// with respect to last pushed loop block
	t.st.Loopend = append(t.st.Loopend, state.LoopEntry{End: endBlock})
	t.ProcessBlock(fn, bodyBlock, st.Body)
	t.st.Loopend = t.st.Loopend[:len(t.st.Loopend)-1]

	if bodyBlock.N.Term == nil {
		bodyBlock.N.NewBr(copyOfCondEntry.N)
	}

	bh.Update(endBlock.V, endBlock.N)
}
