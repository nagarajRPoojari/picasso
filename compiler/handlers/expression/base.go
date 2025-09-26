package expression

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type ExpressionHandler struct {
	st *state.State
}

var ExpressionHandlerInst *ExpressionHandler

func NewExpressionHandler(st *state.State) *ExpressionHandler {
	initOpLookUpTables()
	return &ExpressionHandler{
		st: st,
	}
}

// processExpression handles binary expressions, function calls, member operations etc..
func (t *ExpressionHandler) ProcessExpression(block *ir.Block, expI ast.Expression) (tf.Var, *ir.Block) {
	if expI == nil {
		return tf.NewNullVar(types.NewPointer(types.NewStruct())), block
	}

	switch ex := expI.(type) {

	case ast.SymbolExpression:
		// search for variable locally then gloablly
		return t.processSymbolExpression(ex), block

	case ast.ListExpression:
		// should handle, [[1,2,3], [4,5,6]]

	case ast.NumberExpression:
		// produce a runtime mutable var for the literal (double)
		// by default number will be wrapped up with float64
		return t.st.TypeHandler.BuildVar(block, tf.FLOAT64, constant.NewFloat(types.Double, ex.Value)), block

	case ast.StringExpression:
		formatStr := ex.Value
		strConst := constant.NewCharArrayFromString(formatStr + "\x00")
		global := t.st.Module.NewGlobalDef("", strConst)

		gep := block.NewGetElementPtr(
			global.ContentType,
			global,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0),
		)

		return tf.NewString(block, gep), block

	case ast.NewExpression:
		return t.ProcessNewExpression(block, ex)

	case ast.MemberExpression:
		return t.ProcessMemberExpression(block, ex), block

	case ast.ComputedExpression:
		// e.g, arr[1], arr[1][a.id()];

	case ast.PrefixExpression:
		return t.ProcessPrefixExpression(block, ex), block

	case ast.CallExpression:
		return t.CallFunc(block, ex)

	case ast.BinaryExpression:
		return t.ProcessBinaryExpression(block, ex)
	}

	panic("error")
}
