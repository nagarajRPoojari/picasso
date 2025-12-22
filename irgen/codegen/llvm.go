package generator

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/expression"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	interfaceh "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/interface"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/statement"
	rterr "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/private/runtime"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/pipeline"
)

const (
	TARGETARCH   = "aarch64-unknown-linux-gnu"
	TARGETLAYOUT = "e-m:e-i64:64-n32:64-S128"
)

// LLVM parses abstract syntax tree to generate llvm IR
type LLVM struct {
	st         *state.State
	ModuleName string
}

func NewLLVM(pkgName string, outputDir string) *LLVM {
	m := ir.NewModule()
	m.SourceFilename = pkgName

	rterr.NewErrorHandler(m)
	c.NewInterface(m)

	st := state.NewCompileState(outputDir, pkgName, m)

	expression.ExpressionHandlerInst = expression.NewExpressionHandler(st)
	statement.StatementHandlerInst = statement.NewStatementHandler(st)
	funcs.FuncHandlerInst = funcs.NewFuncHandler(st)
	block.BlockHandlerInst = block.NewBlockHandler(st)
	class.ClassHandlerInst = class.NewClassHandler(st)
	interfaceh.InterfaceHandlerInst = interfaceh.NewInterfaceHandler(st)

	m.TargetTriple = TARGETARCH
	m.DataLayout = TARGETLAYOUT

	return &LLVM{st: st, ModuleName: pkgName}
}

func (t *LLVM) AddImportEntry(entry state.PackageEntry) {
	t.st.Imports[entry.Alias] = entry
}

func (t *LLVM) Dump(outputDir string, file string) {
	f, err := os.Create(fmt.Sprintf("%s/%s.ll", outputDir, file))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(t.st.Module.String())
}

func (t *LLVM) ParseAST(tree ast.BlockStatement) {
	pipeline.NewPipeline(t.st, tree).Run(state.PackageEntry{Name: t.ModuleName, Alias: t.ModuleName})
}
