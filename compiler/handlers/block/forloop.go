package block

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *BlockHandler) processForBlock(fn *ir.Func, bh tf.BlockHolder, st *ast.ForeachStatement) tf.BlockHolder {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	lowerExpr := st.Iterable.(ast.RangeExpression).Lower
	upperExpr := st.Iterable.(ast.RangeExpression).Upper

	indexVal, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, lowerExpr)
	bh = safe

	casted, safeN := t.st.TypeHandler.ImplicitTypeCast(bh.N, tf.INT, indexVal.Load(bh.N))
	bh.N = safeN
	indexVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	iPtr := indexVal.Slot()
	t.st.Vars.AddNewVar(st.Value, indexVal)

	upperVal, safe := expression.ExpressionHandlerInst.ProcessExpression(bh, upperExpr)
	bh = safe
	casted, safeN = t.st.TypeHandler.ImplicitTypeCast(bh.N, tf.INT, upperVal.Load(bh.N))
	bh.N = safeN
	upperVal = t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.INT), casted)

	loopCond := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}
	loopBody := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}
	loopInc := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}
	loopEnd := tf.BlockHolder{V: bh.V, N: fn.NewBlock("")}

	bh.N.NewBr(loopCond.N)

	iVal := loopCond.N.NewLoad(types.I64, iPtr)
	cond := loopCond.N.NewICmp(enum.IPredSLT, iVal, upperVal.Load(loopCond.N))
	loopCond.N.NewCondBr(cond, loopBody.N, loopEnd.N)

	bodyBlock := t.ProcessBlock(fn, loopBody, st.Body)
	bodyBlock.N.NewBr(loopInc.N)

	iVal2 := loopInc.N.NewLoad(types.I64, iPtr)
	iNext := loopInc.N.NewAdd(iVal2, constant.NewInt(types.I64, 1))
	loopInc.N.NewStore(iNext, iPtr)
	loopInc.N.NewBr(loopCond.N)

	return loopEnd
}
