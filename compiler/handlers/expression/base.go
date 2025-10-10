package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
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

// ProcessExpression evaluates an AST expression node within the given IR block
// and produces a corresponding runtime variable along with the (possibly updated) block.
//
// It serves as a central dispatcher, delegating handling of each expression type
// (symbols, literals, function calls, object creation, operators, etc.) to the
// appropriate specialized method.
//
// Parameters:
//
//	block - the current IR block in which code generation should occur
//	expI  - the AST expression to evaluate (nil-safe)
//
// Returns:
//
//	tf.Var     - the resulting typed variable representing the expression
//	*ir.Block  - the (possibly modified) IR block after processing
func (t *ExpressionHandler) ProcessExpression(bh tf.BlockHolder, expI ast.Expression) (tf.Var, tf.BlockHolder) {
	if expI == nil {
		return tf.NewNullVar(types.NewPointer(types.NewStruct())), bh
	}

	switch ex := expI.(type) {

	case ast.SymbolExpression:
		if t.st.TypeHandler.Exists(ex.Value) {
			return t.st.TypeHandler.BuildVar(bh, tf.NewType(ex.Value), nil), bh
		}
		return t.processSymbolExpression(ex), bh

	case ast.ListExpression:
		// should handle, [[1,2,3], [4,5,6]]
		// @todo

	case ast.NumberExpression:
		// by default number will be wrapped up with float64
		return t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.FLOAT64), constant.NewFloat(types.Double, ex.Value)), bh

	case ast.StringExpression:
		return t.ProcessStringLiteral(bh, ex), bh

	case ast.NewExpression:
		return t.ProcessNewExpression(bh, ex)

	case ast.MemberExpression:
		return t.ProcessMemberExpression(bh, ex)

	case ast.ComputedExpression:
		return t.ProcessIndexingExpression(bh, ex)

	case ast.PrefixExpression:
		return t.ProcessPrefixExpression(bh, ex), bh

	case ast.CallExpression:
		return t.CallFunc(bh, ex)

	case ast.BinaryExpression:
		return t.ProcessBinaryExpression(bh, ex)
	}

	errorutils.Abort(errorutils.InvalidExpression)
	return nil, tf.BlockHolder{}
}
