package generator

import (
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/scope"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	rterr "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/private/runtime"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/pipeline"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
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
