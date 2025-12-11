package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

// ProcessStringLiteral creates a runtime string variable from a string literal.
//
// Parameters:
//
//	block - the current IR block
//	ex    - AST StringExpression node
//
// Returns:
//
//	tf.Var - runtime string variable
func (t *ExpressionHandler) ProcessStringLiteral(bh *bc.BlockHolder, ex ast.StringExpression) tf.Var {
	formatStr := ex.Value

	strConst := constant.NewCharArrayFromString(formatStr + "\x00")
	global := t.st.Module.NewGlobalDef("", strConst)
	gep := bh.N.NewGetElementPtr(
		global.ContentType,
		global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	malloc := t.st.CI.Funcs[c.FUNC_ALLOC]
	size := constant.NewInt(types.I64, int64(len(formatStr)+1))
	heapPtr := bh.N.NewCall(malloc, size)

	memcpy := t.st.CI.Funcs[c.FUNC_MEMCPY]
	i8ptr := types.NewPointer(types.I8)
	src := bh.N.NewBitCast(gep, i8ptr)
	dest := bh.N.NewBitCast(heapPtr, i8ptr)
	bh.N.NewCall(memcpy, dest, src, size, constant.NewBool(false))

	return tf.NewString(bh, heapPtr)
}
