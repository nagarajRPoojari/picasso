package block

import (
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// ProcessBlock serves as the central recursive dispatcher for transforming a sequence
// of AST statements into LLVM IR instructions.
//
// Lifecycle:
//  1. Scope Management: Automatically manages lexical scoping by pushing a new
//     symbol table block on entry and popping it via 'defer' on exit.
//  2. Dispatching: Iterates through the AST node slice, delegating specific
//     codegen logic to specialized handlers (Statement, Expression, etc.)
//     via the Mediator.
//  3. Control Flow: Handles early termination for 'Break' and 'Return'
//     statements to ensure subsequent unreachable AST nodes are not
//     translated into the current basic block.
func (t *BlockHandler) ProcessBlock(fn *ir.Func, bh *bc.BlockHolder, sts []ast.Statement) {
	t.st.Vars.AddBlock()
	defer t.st.Vars.RemoveBlock()

	// Cache the handler to avoid repeated type assertions in the loop
	sh := t.m.GetStatementHandler().(*statement.StatementHandler)

	for _, stI := range sts {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			sh.DeclareVariable(bh, &st)

		case ast.ExpressionStatement:
			t.processExpressionStatement(st, bh, sh)

		case ast.IfStatement:
			t.processIfElseBlock(fn, bh, &st, nil)

		case ast.ForeachStatement:
			t.processForBlock(fn, bh, &st)

		case ast.WhileStatement:
			t.processWhileBlock(fn, bh, &st)

		case ast.BreakStatement:
			t.handleBreak(bh)
			return // Stop processing this block after a break

		case ast.ReturnStatement:
			sh.Return(bh, &st, t.getRetType(fn))
			return // Stop processing this block after a return
		}
	}
}

// processExpressionStatement handles standalone expressions (e.g, func calls, variable assignments etc.)
func (t *BlockHandler) processExpressionStatement(st ast.ExpressionStatement, bh *bc.BlockHolder, sh *statement.StatementHandler) {
	switch exp := st.Expression.(type) {
	case ast.AssignmentExpression:
		sh.AssignVariable(bh, &exp)
	case ast.CallExpression:
		sh.CallFunc(bh, exp)
	case ast.NewExpression:
		sh.ProcessNewExpression(bh, exp)
	default:
		errorutils.Abort(errorutils.InvalidStatement)
	}
}

// handleBreak verifies if there is any breakable block (aka loop block)
// if so, return to loopend block else throw InvalidBreakStatement error
func (t *BlockHandler) handleBreak(bh *bc.BlockHolder) {
	if len(t.st.Loopend) == 0 {
		errorutils.Abort(errorutils.InvalidBreakStatement)
	}

	loopend := t.st.Loopend[len(t.st.Loopend)-1]
	bh.N.NewBr(loopend.End.N)
}

// get return type of function based on explicit store in compiler state
// Note: can't depend on fn.Sig.RetType since i'll lose extra informations
// like signed/unsigned status.
func (t *BlockHandler) getRetType(fn *ir.Func) ast.Type {
	name := fn.Name()
	idx := strings.LastIndex(name, ".")
	if idx == -1 {
		return nil // Or handle global functions
	}

	clsName := name[:idx]
	if clsMeta, ok := t.st.Classes[clsName]; ok {
		return clsMeta.Returns[name]
	}
	return nil
}
