package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-ini/ini"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
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
			if !stc.IsBasePkg() {
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
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.IsBasePkg() {
				t.importBasePackages(llvm.st.LibMethods, stc.EndName())
			}
			llvm.AddImportEntry(state.ImportEntry{Name: stc.Name, Identifier: stc.Alias})
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
		if stc, ok := st.(ast.ImportStatement); ok && stc.IsBasePkg() {
			t.importBasePackages(llvm.st.LibMethods, stc.EndName())
		}
	}

	pipeline.NewPipeline(llvm.st, pkgAST).Declare(pkgName)
}

// importBasePackages resolve base module imports.
func (t *generator) importBasePackages(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, module)
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("builtin.%s.%s", module, name)
		methodMap[n] = f
	}
}

// compile passes tree through all steps to output IR
func (t *generator) compile(tree ast.BlockStatement, llvm *LLVM) {
	llvm.ParseAST(tree)
}

// LoadPackages loads package with their AST by going through project.ini file
func LoadPackages(projectIniPath string) map[string]ast.BlockStatement {
	cfg, err := ini.Load(projectIniPath)
	if err != nil {
		panic("failed to read project.ini")
	}

	rootDir := cfg.Section("root").Key("path").String()
	if rootDir == "" {
		panic("root.path is empty in project.ini")
	}

	pkgs := make(map[string]ast.BlockStatement)

	err = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".pic" {
			return nil
		}

		sourceBytes, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("unable to read %s: %w", path, err)
		}

		source := string(sourceBytes)
		tree := parser.Parse(source)

		// relative path from root
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// remove extension
		rel = strings.TrimSuffix(rel, ".pic")

		// normalize to forward slashes for package name
		pkgName := filepath.ToSlash(rel)

		pkgs[pkgName] = tree
		return nil
	})

	if err != nil {
		panic(err)
	}

	return pkgs
}
