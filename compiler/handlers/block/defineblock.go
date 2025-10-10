package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/statement"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// ProcessBlock builds LLVM IR for a list of AST statements within a function.
// It creates a new variable scope, processes each statement in order, and
// emits the corresponding IR instructions (variable declarations, expressions,
// conditionals, and returns).
//
// Params:
//
//	fn    – the LLVM function being built
//	entry – the current IR basic block where instructions are inserted
//	sts   – slice of AST statements to process
//
// Return:
//
//	The updated IR basic block after processing all statements.
func (t *BlockHandler) ProcessBlock(fn *ir.Func, bh tf.BlockHolder, sts []ast.Statement) tf.BlockHolder {
	// add new scope block for variables
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	for _, stI := range sts {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			bh = statement.StatementHandlerInst.DeclareVariable(bh, &st)

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				bh = statement.StatementHandlerInst.AssignVariable(bh, &exp)
			case ast.CallExpression:
				_, bh = statement.StatementHandlerInst.CallFunc(bh, exp)
			case ast.NewExpression:
				_, bh = statement.StatementHandlerInst.ProcessNewExpression(bh, exp)
			default:
				errorutils.Abort(errorutils.InvalidStatement)
			}

		case ast.IfStatement:
			bh = t.processIfElseBlock(fn, bh, &st)
		case ast.ForeachStatement:
			bh = t.processForBlock(fn, bh, &st)
		case ast.ReturnStatement:
			retType := fn.Sig.RetType
			statement.StatementHandlerInst.Return(bh, &st, retType)
		}
	}

	return bh
}
