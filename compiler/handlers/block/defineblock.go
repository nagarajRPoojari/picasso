package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/statement"
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
func (t *BlockHandler) ProcessBlock(fn *ir.Func, entry *ir.Block, sts []ast.Statement) *ir.Block {
	// add new scope block for variables
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	for _, stI := range sts {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			entry = statement.StatementHandlerInst.DeclareVariable(entry, &st)

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				entry = statement.StatementHandlerInst.AssignVariable(entry, &exp)
			case ast.CallExpression:
				_, entry = statement.StatementHandlerInst.CallFunc(entry, exp)
			case ast.NewExpression:
				_, entry = statement.StatementHandlerInst.ProcessNewExpression(entry, exp)
			default:
				errorutils.Abort(errorutils.InvalidStatement)
			}

		case ast.IfStatement:
			entry = t.processIfElseBlock(fn, entry, &st)
		case ast.ReturnStatement:
			retType := fn.Sig.RetType
			statement.StatementHandlerInst.Return(entry, &st, retType)
		}
	}

	return entry
}
