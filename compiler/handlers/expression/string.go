package expression

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
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
func (t *ExpressionHandler) ProcessStringLiteral(block *ir.Block, ex ast.StringExpression) tf.Var {
	formatStr := ex.Value
	strConst := constant.NewCharArrayFromString(formatStr + "\x00")
	global := t.st.Module.NewGlobalDef("", strConst)

	gep := block.NewGetElementPtr(
		global.ContentType,
		global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	return tf.NewString(block, gep)
}
