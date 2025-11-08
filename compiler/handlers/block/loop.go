package block

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

func (t *BlockHandler) processForBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.ForeachStatement) {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	lowerExpr := st.Iterable.(ast.RangeExpression).Lower
	upperExpr := st.Iterable.(ast.RangeExpression).Upper

	indexVal := expression.ExpressionHandlerInst.ProcessExpression(bh, lowerExpr)

	casted := t.st.TypeHandler.ImplicitTypeCast(bh, tf.INT, indexVal.Load(bh))
	indexVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	iPtr := indexVal.Slot()
	t.st.Vars.AddNewVar(st.Value, indexVal)

	upperVal := expression.ExpressionHandlerInst.ProcessExpression(bh, upperExpr)
	casted = t.st.TypeHandler.ImplicitTypeCast(bh, tf.INT, upperVal.Load(bh))
	upperVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	loopCond := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopBody := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopInc := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	loopEnd := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	bh.N.NewBr(loopCond.N)

	iVal := loopCond.N.NewLoad(types.I64, iPtr)
	cond := loopCond.N.NewICmp(enum.IPredSLT, iVal, upperVal.Load(loopCond))
	loopCond.N.NewCondBr(cond, loopBody.N, loopEnd.N)

	t.st.Loopend = append(t.st.Loopend, state.LoopEntry{End: loopEnd})
	t.ProcessBlock(fn, loopBody, st.Body)
	t.st.Loopend = t.st.Loopend[:len(t.st.Loopend)-1]

	loopBody.N.NewBr(loopInc.N)

	iVal2 := loopInc.N.NewLoad(types.I64, iPtr)
	iNext := loopInc.N.NewAdd(iVal2, constant.NewInt(types.I64, 1))
	loopInc.N.NewStore(iNext, iPtr)
	loopInc.N.NewBr(loopCond.N)

	bh.Update(loopEnd.V, loopEnd.N)
}

func (t *BlockHandler) processWhileBlock(fn *ir.Func, bh *bc.BlockHolder, st *ast.WhileStatement) {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	condBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	bodyBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))
	endBlock := bc.NewBlockHolder(bh.V, fn.NewBlock(""))

	bh.N.NewBr(condBlock.N)

	res := expression.ExpressionHandlerInst.ProcessExpression(condBlock, st.Condition)
	cond := res.Load(condBlock)

	condBlock.N.NewCondBr(cond, bodyBlock.N, endBlock.N)

	t.st.Loopend = append(t.st.Loopend, state.LoopEntry{End: endBlock})
	t.ProcessBlock(fn, bodyBlock, st.Body)
	t.st.Loopend = t.st.Loopend[:len(t.st.Loopend)-1]

	bodyBlock.N.NewBr(condBlock.N)

	bh.Update(endBlock.V, endBlock.N)
}
