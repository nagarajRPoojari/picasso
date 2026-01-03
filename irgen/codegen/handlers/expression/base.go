// Package expression provides the logic for evaluating AST expressions and
// lowering them into LLVM IR values. It serves as the primary recursive
// engine for the backend, translating high-level constructs like function
// calls, member access, and arithmetic operations into addressable
// memory slots or register values.
package expression

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/contract"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// ExpressionHandler encapsulates the state required to generate IR
// for diverse expressions. It maintains a reference to the global
// compiler state to resolve types, symbols, and class metadata.
type ExpressionHandler struct {
	st *state.State
	m  contract.Mediator
}

// NewExpressionHandler initializes the handler and populates the
// operator lookup tables (arithmetic, logical, and comparison)
// used during binary expression processing.
func NewExpressionHandler(st *state.State, m contract.Mediator) *ExpressionHandler {
	initOpLookUpTables()
	return &ExpressionHandler{
		st: st,
		m:  m,
	}
}

// ProcessExpression acts as the central dispatcher for the expression
// sub-system. It performs a type switch on the AST node to delegate
// code generation to specialized handlers, returning a tf.Var
// which abstracts the underlying LLVM value and its Niyama type.
//
// Technical Logic:
//   - Recursion: Evaluates complex, nested expressions by drilling
//     down into sub-expressions (e.g., arguments in a call).
//   - Null Handling: Automatically wraps nil or null expressions
//     into an opaque pointer structure.
//   - Default Numerics: Coerces raw number literals to float64 (Double)
//     to maintain the language's uniform numeric behavior.
func (t *ExpressionHandler) ProcessExpression(bh *bc.BlockHolder, expI ast.Expression) tf.Var {
	if expI == nil {
		return tf.NewNullVar(types.NewPointer(types.NewStruct()))
	}

	switch ex := expI.(type) {

	case ast.NullExpression:
		return tf.NewNullVar(types.NewPointer(types.NewStruct()))

	case ast.SymbolExpression:
		if ret, ok := t.loopUpTypeTable(bh, ex.Value); ok {
			return ret
		}
		return t.processSymbolExpression(ex)

	case ast.ListExpression:
		// @todo should handle, [[1,2,3], [4,5,6]]

	case ast.NumberExpression:
		// by default number will be wrapped up with float64
		return t.st.TypeHandler.BuildVar(bh, tf.NewType(tf.FLOAT64), constant.NewFloat(types.Double, ex.Value))

	case ast.StringExpression:
		return t.ProcessStringLiteral(bh, ex)

	case ast.NewExpression:
		return t.ProcessNewExpression(bh, ex)

	case ast.MemberExpression:
		if m, ok := ex.Member.(ast.SymbolExpression); ok {
			if ret, ok := t.loopUpTypeTable(bh, fmt.Sprintf("%s.%s", m.Value, ex.Property)); ok {
				return ret
			}
		}
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

func (t *ExpressionHandler) loopUpTypeTable(bh *bc.BlockHolder, val string) (tf.Var, bool) {
	if t.st.TypeHandler.Exists(val) {
		return t.st.TypeHandler.BuildVar(bh, tf.NewType(val), nil), true
	}
	return nil, false
}
