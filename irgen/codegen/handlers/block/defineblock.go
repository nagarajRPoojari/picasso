package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
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
func (t *BlockHandler) ProcessBlock(fn *ir.Func, bh *bc.BlockHolder, sts []ast.Statement) {
	// add new scope block for variables
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	for _, stI := range sts {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			statement.StatementHandlerInst.DeclareVariable(bh, &st)

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				statement.StatementHandlerInst.AssignVariable(bh, &exp)
			case ast.CallExpression:
				statement.StatementHandlerInst.CallFunc(bh, exp)
			case ast.NewExpression:
				statement.StatementHandlerInst.ProcessNewExpression(bh, exp)
			default:
				errorutils.Abort(errorutils.InvalidStatement)
			}

		case ast.IfStatement:
			t.processIfElseBlock(fn, bh, &st)
		case ast.ForeachStatement:
			t.processForBlock(fn, bh, &st)
		case ast.WhileStatement:
			t.processWhileBlock(fn, bh, &st)
		case ast.BreakStatement:
			if len(t.st.Loopend) == 0 {
				errorutils.Abort(errorutils.InvalidBreakStatement)
			}
			loopend := t.st.Loopend[len(t.st.Loopend)-1]
			bh.N.NewBr(loopend.End.N)
			// bh.Update(loopend.End.V, loopend.End.N)

			// return from this block immediately, simply ignoring all upcomming statements
			return
		case ast.ReturnStatement:
			retType := fn.Sig.RetType
			statement.StatementHandlerInst.Return(bh, &st, retType)

			// return from this block immediately, simply ignoring all upcomming statements
			return
		}
	}
}
