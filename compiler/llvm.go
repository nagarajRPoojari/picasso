package compiler

import (
	"os"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/gc"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/block"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/class"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/expression"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/identifier"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/scope"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/statement"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	rterr "github.com/nagarajRPoojari/x-lang/compiler/libs/private/runtime"
	"github.com/nagarajRPoojari/x-lang/compiler/pipeline"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// LLVM parses abstract syntax tree to generate llvm IR
type LLVM struct {
	st *state.State
}

func NewLLVM() *LLVM {
	m := ir.NewModule()
	tree := scope.NewVarTree()
	// tree.()

	st := &state.State{
		Module:            m,
		TypeHandler:       tf.NewTypeHandler(),
		Vars:              tree,
		Methods:           make(map[string]*ir.Func),
		Classes:           make(map[string]*tf.MetaClass),
		IdentifierBuilder: identifier.NewIdentifierBuilder(MAIN),
		LibMethods:        make(map[string]function.Func),
		GC:                gc.GetGC(m),
	}

	expression.ExpressionHandlerInst = expression.NewExpressionHandler(st)
	statement.StatementHandlerInst = statement.NewStatementHandler(st)
	funcs.FuncHandlerInst = funcs.NewFuncHandler(st)
	block.BlockHandlerInst = block.NewBlockHandler(st)
	class.ClassHandlerInst = class.NewClassHandler(st)

	rterr.Instance = rterr.NewErrorHandler(m)

	return &LLVM{st: st}
}

func (t *LLVM) Dump(file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(t.st.Module.String())
}

func (t *LLVM) ParseAST(tree ast.BlockStatement) {
	pipeline.NewPipeline(t.st, tree).Run()
}
