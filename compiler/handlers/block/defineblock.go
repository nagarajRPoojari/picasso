package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/statement"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *BlockHandler) ProcessBlock(fn *ir.Func, entry *ir.Block, sts []ast.Statement) *ir.Block {
	// add new block
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	for _, stI := range sts {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			statement.StatementHandlerInst.DeclareVariable(entry, &st)

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				statement.StatementHandlerInst.AssignVariable(entry, &exp)
			case ast.CallExpression:
				statement.StatementHandlerInst.CallFunc(entry, exp)
			case ast.NewExpression:
				statement.StatementHandlerInst.ProcessNewExpression(entry, exp)
			default:
				errorsx.PanicCompilationError("invalid expression statement")
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
