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
	llvms     map[string]*LLVM
	outputDir string

	// builtPkgs tracks packages that are completely built and compiled
	builtPkgs map[string]struct{}
	// visitingPkgs tracks packages currently in the recursion stack (for cycle detection)
	visitingPkgs map[string]struct{}
}

func NewGenerator(path string, outputDir string) *generator {
	return &generator{
		packages:     LoadPackages(path),
		llvms:        make(map[string]*LLVM),
		outputDir:    outputDir,
		builtPkgs:    make(map[string]struct{}),
		visitingPkgs: make(map[string]struct{}),
	}
}

func (t *generator) BuildAll() {
	t.buildPackage("start")
}

// buildPackage implements Post-Order Traversal (Bottom-Up) to ensure global state integrity.
func (t *generator) buildPackage(pkgName string) {
	if _, ok := t.builtPkgs[pkgName]; ok {
		return
	}

	if _, ok := t.visitingPkgs[pkgName]; ok {
		panic(fmt.Sprintf("Cyclic dependency detected involving package: %s", pkgName))
	}
	t.visitingPkgs[pkgName] = struct{}{}

	tree, ok := t.packages[pkgName]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, pkgName)
		return
	}

	directUserImports := t.extractUserImports(tree)
	for _, imp := range directUserImports {
		t.buildPackage(imp)
	}

	// Create new LLVM context for this package (Safe, as children are finished)
	llvm := NewLLVM(pkgName)
	t.llvms[pkgName] = llvm

	// Resolve Imports: Declare symbols from direct and transitive dependencies (B and C)
	t.resolveImports(tree, directUserImports, llvm)

	// Compile
	t.compile(tree, llvm)

	// Dump
	llvm.Dump(t.outputDir, pkgName)

	// Mark as finished
	delete(t.visitingPkgs, pkgName)
	t.builtPkgs[pkgName] = struct{}{}
}

// extractUserImports returns only the names of non-builtin imported packages.
func (t *generator) extractUserImports(tree ast.BlockStatement) []string {
	var imports []string
	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.From != constants.BUILTIN {
				imports = append(imports, stc.Name)
			}
		}
	}
	return imports
}

func (t *generator) resolveImports(tree ast.BlockStatement, directUserImports []string, llvm *LLVM) {
	// declared map tracks all packages added to this module's symbol table to prevent redundancy
	declared := make(map[string]struct{})

	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok && stc.From == constants.BUILTIN {
			t.importBasePackages(llvm.st.LibMethods, stc.Name)
		}
	}

	for _, pkgName := range directUserImports {
		t.recursiveTransitiveDeclaration(pkgName, llvm, declared)
	}
}

// recursiveTransitiveDeclaration declares symbols of pkgName and all its dependencies (C)
// into the current module (A). This is the key to fixing the transitive dependency issue.
func (t *generator) recursiveTransitiveDeclaration(pkgName string, llvm *LLVM, declared map[string]struct{}) {
	if _, ok := declared[pkgName]; ok {
		return
	}
	declared[pkgName] = struct{}{}

	pkgAST := t.packages[pkgName]

	subImports := t.extractUserImports(pkgAST)
	for _, sub := range subImports {
		t.recursiveTransitiveDeclaration(sub, llvm, declared)
	}

	for _, st := range pkgAST.Body {
		if stc, ok := st.(ast.ImportStatement); ok && stc.From == constants.BUILTIN {
			t.importBasePackages(llvm.st.LibMethods, stc.Name)
		}
	}

	pipeline.NewPipeline(llvm.st, pkgAST).Declare()
}

func (t *generator) importBasePackages(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, module)
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s", module, name)
		methodMap[n] = f
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
