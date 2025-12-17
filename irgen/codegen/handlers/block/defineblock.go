package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// ProcessBlock translates a sequence of AST statements into LLVM IR instructions.
// It acts as the primary dispatcher for the code generator, delegating specific
// statement types to their respective handlers while managing the lifecycle
// of lexical scopes and control flow integrity.
//
// Key logic:
//   - Scope Management: Pushes a new variable scope onto the symbol table stack
//     on entry and ensures its removal via 'defer' to prevent memory leaks or
//     variable shadowing issues.
//   - Dispatch Logic: Switches on concrete AST types to handle declarations,
//     assignments, function calls, and complex control structures (If, While, For).
//   - Early Termination: Detects 'Break' and 'Return' statements to immediately
//     halt instruction emission for the current block, preventing unreachable
//     code errors in the generated IR.
//   - Loop Context: Validates 'Break' statements against the Loopend stack to
//     ensure they only occur within valid iterative contexts.
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
