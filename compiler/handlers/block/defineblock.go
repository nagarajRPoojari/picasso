package block

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
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
			expression.ExpressionHandlerInst.DeclareVariable(entry, &st)

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				expression.ExpressionHandlerInst.AssignVariable(entry, &exp)
			case ast.CallExpression:
				expression.ExpressionHandlerInst.CallFunc(entry, exp)
			default:
				errorsx.PanicCompilationError("invalid statement")
			}

		case ast.IfStatement:
			entry = t.processIfElseStatement(fn, entry, &st)
		case ast.ReturnStatement:
			retType := fn.Sig.RetType
			statement.StatementHandlerInst.Return(entry, &st, retType)
		}
	}

	return entry
}
