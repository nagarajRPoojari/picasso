package generator

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/c"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/contract"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/mediator"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
	rterr "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/private/runtime"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/pipeline"
)

// const (
// 	TARGETARCH   = "aarch64-unknown-linux-gnu"
// 	TARGETLAYOUT = "e-m:e-i64:64-n32:64-S128"
// )

const (
	TARGETARCH   = "arm64-apple-darwin"
	TARGETLAYOUT = "e-m:o-i64:64-i128:128-n32:64-S128"
)

// LLVM parses abstract syntax tree to generate llvm IR
type LLVM struct {
	st         *state.State
	ModuleName string
	m          contract.Mediator
}

func NewLLVM(pkgName string, outputDir string) *LLVM {
	module := ir.NewModule()
	module.SourceFilename = pkgName
	module.TargetTriple = TARGETARCH
	module.DataLayout = TARGETLAYOUT

	rterr.InitErrorHandler(module)
	c.InitInterface(module)

	st := state.NewCompileState(outputDir, pkgName, module)
	m := mediator.InitMediator(st)

	return &LLVM{st: st, m: m, ModuleName: pkgName}
}

func (t *LLVM) GetModule() *ir.Module {
	return t.st.Module
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
	pipeline.NewPipeline(t.st, t.m, tree).Run(state.PackageEntry{Name: t.ModuleName, Alias: t.ModuleName})
}
