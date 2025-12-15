package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	errorutils "github.com/nagarajRPoojari/niyama/frontend/codegen/error"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/state"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
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
func (t *ExpressionHandler) ProcessExpression(bh *bc.BlockHolder, expI ast.Expression) tf.Var {
	if expI == nil {
		return tf.NewNullVar(types.NewPointer(types.NewStruct()))
	}

	switch ex := expI.(type) {

	case ast.NullExpression:
		return tf.NewNullVar(types.NewPointer(types.NewStruct()))

	case ast.SymbolExpression:
		if t.st.TypeHandler.Exists(ex.Value) {
			return t.st.TypeHandler.BuildVar(bh, tf.NewType(ex.Value), nil)
		}
		return t.processSymbolExpression(ex)

	case ast.ListExpression:
		// should handle, [[1,2,3], [4,5,6]]
		// @todo

	case ast.NumberExpression:
		// by default number will be wrapped up with float64
		return t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.FLOAT64), constant.NewFloat(types.Double, ex.Value))

	case ast.StringExpression:
		return t.ProcessStringLiteral(bh, ex)

	case ast.NewExpression:
		return t.ProcessNewExpression(bh, ex)

	case ast.MemberExpression:
		return t.ProcessMemberExpression(bh, ex)

	case ast.ComputedExpression:
		return t.ProcessIndexingExpression(bh, ex)

	case ast.PrefixExpression:
		return t.ProcessPrefixExpression(bh, ex)

	case ast.CallExpression:
		return t.CallFunc(bh, ex)

	case ast.BinaryExpression:
		return t.ProcessBinaryExpression(bh, ex)
	}

	errorutils.Abort(errorutils.InvalidExpression)
	return nil
}
