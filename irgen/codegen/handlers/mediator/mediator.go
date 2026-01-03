package mediator

import (
	"github.com/nagarajRPoojari/niyama/irgen/codegen/contract"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	interfaceh "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/interface"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
)

type Mediator struct {
	expHandler       *expression.ExpressionHandler
	stmtHandler      *statement.StatementHandler
	funcHandler      *funcs.FuncHandler
	blockHandler     *block.BlockHandler
	classHandler     *class.ClassHandler
	interfaceHandler *interfaceh.InterfaceHandler
}

func InitMediator(st *state.State) contract.Mediator {
	m := &Mediator{}

	// Inject the mediator (m) into every handler
	m.expHandler = expression.NewExpressionHandler(st, m)
	m.stmtHandler = statement.NewStatementHandler(st, m)
	m.funcHandler = funcs.NewFuncHandler(st, m)
	m.blockHandler = block.NewBlockHandler(st, m)
	m.classHandler = class.NewClassHandler(st, m)
	m.interfaceHandler = interfaceh.NewInterfaceHandler(st, m)

	return m
}

// Getters implement the contract.Mediator interface
func (m *Mediator) GetExpressionHandler() any { return m.expHandler }
func (m *Mediator) GetStatementHandler() any  { return m.stmtHandler }
func (m *Mediator) GetFuncHandler() any       { return m.funcHandler }
func (m *Mediator) GetBlockHandler() any      { return m.blockHandler }
func (m *Mediator) GetClassHandler() any      { return m.classHandler }
func (m *Mediator) GetInterfaceHandler() any  { return m.interfaceHandler }
