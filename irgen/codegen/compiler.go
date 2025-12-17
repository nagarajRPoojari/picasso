package generator

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/pipeline"
	"github.com/nagarajRPoojari/niyama/irgen/parser"
)

type generator struct {
	packages  map[string]ast.BlockStatement
	outputDir string
}

func NewGenerator(path string, outputDir string) *generator {
	return &generator{packages: LoadPackages(path), outputDir: outputDir}
}

func (t *generator) BuildAll() {
	for name, pkg := range t.packages {
		llvm := NewLLVM(name)
		t.build(pkg, llvm)
		t.compile(pkg, llvm)

		llvm.Dump(t.outputDir, name)
	}
}

func (t *generator) build(tree ast.BlockStatement, llvm *LLVM) {
	userModules := map[string]struct{}{}
	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.From == constants.BUILTIN {
				t.importBaseModules(llvm.st.LibMethods, stc.Name)
			} else {
				userModules[stc.Name] = struct{}{}
			}
		}
	}
	t.ImportUserModules(userModules, llvm)
}

func (t *generator) importBaseModules(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, module)
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s", module, name)
		methodMap[n] = f
	}
}

func (t *generator) ImportUserModules(userModules map[string]struct{}, llvm *LLVM) {
	for pkgName, pkg := range t.packages {
		if _, ok := userModules[pkgName]; !ok {
			continue
		}
		pipeline.NewPipeline(llvm.st, pkg).Declare()
	}
}

func (t *generator) compile(tree ast.BlockStatement, llvm *LLVM) {
	llvm.ParseAST(tree)
}

func LoadPackages(packagePath string) map[string]ast.BlockStatement {
	cfg, err := ini.Load(packagePath)
	if err != nil {
		panic("failed to read package.ini")
	}

	m := make(map[string]ast.BlockStatement, 0)

	for _, section := range cfg.Sections()[1:] {
		sourceBytes, err := os.ReadFile(section.Key("path").String())
		if err != nil {
			panic(fmt.Sprintf("unable to read: %v", err))
		}
		source := string(sourceBytes)
		tree := parser.Parse(source)
		m[section.Name()] = tree
	}

	return m
}
