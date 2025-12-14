package generator

import (
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/block"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/class"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/expression"
	funcs "github.com/nagarajRPoojari/x-lang/generator/handlers/func"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/identifier"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/scope"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/state"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/statement"
	function "github.com/nagarajRPoojari/x-lang/generator/libs/func"
	rterr "github.com/nagarajRPoojari/x-lang/generator/libs/private/runtime"
	"github.com/nagarajRPoojari/x-lang/generator/pipeline"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
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
		Classes:           make(map[string]*tf.MetaClass),
		IdentifierBuilder: identifier.NewIdentifierBuilder(MAIN),
		LibMethods:        make(map[string]function.Func),
		CI:                c.NewInterface(m),
	}

	expression.ExpressionHandlerInst = expression.NewExpressionHandler(st)
	statement.StatementHandlerInst = statement.NewStatementHandler(st)
	funcs.FuncHandlerInst = funcs.NewFuncHandler(st)
	block.BlockHandlerInst = block.NewBlockHandler(st)
	class.ClassHandlerInst = class.NewClassHandler(st)

	rterr.Instance = rterr.NewErrorHandler(m)

	m.TargetTriple = "aarch64-unknown-linux-gnu"
	m.DataLayout = "e-m:e-i64:64-n32:64-S128"

	// @todo: this is not the right place
	ax := types.NewStruct(
		types.NewPointer(types.I8),  // data
		types.NewPointer(types.I64), // shape pointer
		types.I64,                   // length
		types.I64,                   // rank
	)
	st.Module.NewTypeDef("array", ax)

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
