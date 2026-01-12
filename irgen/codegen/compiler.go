package generator

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/llir/llvm/asm"
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/pipeline"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/tools"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/utils"
	"github.com/nagarajRPoojari/niyama/irgen/parser"
)

type generator struct {
	// ast of all modified packages involved in the project.s
	packages map[string]ast.BlockStatement

	allPkgs map[string]struct{}

	// llvm instance of all packages.
	llvms map[string]*LLVM

	ffiModules map[string]*ir.Module

	// outputDir is directory where IR & info files will be dumped.
	outputDir string

	// cached pkgs
	cachedPkgs map[string]struct{}

	// builtPkgs tracks packages that are completely built and compiled
	// key should be fully qualified name of package. e.g, os.io
	builtPkgs map[string]struct{}

	// visitingPkgs tracks packages currently in the recursion stack (for cycle detection)
	// key should be fully qualified name of package. e.g, os.io
	visitingPkgs map[string]struct{}
}

func NewGenerator(projectDir string) *generator {
	outputDir := filepath.Join(projectDir, BUILD)

	modifiedPkgs, allPkgs := LoadPackages(projectDir)
	return &generator{
		packages:     modifiedPkgs,
		allPkgs:      allPkgs,
		llvms:        make(map[string]*LLVM),
		ffiModules:   make(map[string]*ir.Module),
		outputDir:    outputDir,
		builtPkgs:    make(map[string]struct{}),
		visitingPkgs: make(map[string]struct{}),
	}
}

func (t *generator) BuildAll() {
	// main file is expected to be named as start.pic.
	// @todo: main.pic would be a good choise, why did I even replace
	// all 'main' with 'start'?
	t.buildPackage(state.PackageEntry{Name: MAIN, Alias: MAIN})

	// for all modified packages, generate .exports
	for pkgName := range t.packages {
		t.generateExports(pkgName)
	}
}

func (t *generator) generateExports(pkgName string) {
	outputPath := filepath.Join(t.outputDir, fmt.Sprintf("%s.exports", pkgName))

	// clear definitions in AST
	for _, stmt := range t.packages[pkgName].Body {
		if cls, ok := stmt.(ast.ClassDeclarationStatement); ok {
			for i, stmt := range cls.Body {
				if funcStmt, ok := stmt.(ast.FunctionDefinitionStatement); ok {
					cls.Body[i] = ast.FunctionDefinitionStatement{
						Parameters: funcStmt.Parameters,
						Name:       funcStmt.Name,
						Body:       []ast.Statement{},
						Hash:       funcStmt.Hash,
						ReturnType: funcStmt.ReturnType,
						IsStatic:   funcStmt.IsStatic,
					}
				}
			}
		}
	}

	utils.SaveToFile(outputPath, t.packages[pkgName])
}

// buildPackage implements Post-Order Traversal (Bottom-Up) to ensure global state integrity.
// LLVM maintains some global state unique to that module, traversing bottom up avoids
// such global val overriding.
// maintaining global is a design limitation that need to be fixed in future @todo.
func (t *generator) buildPackage(pkg state.PackageEntry) {
	pkgName := pkg.Name
	if _, ok := t.builtPkgs[pkgName]; ok {
		return
	}

	if _, ok := t.visitingPkgs[pkgName]; ok {
		panic(fmt.Sprintf("Cyclic dependency detected involving package: %s", pkgName))
	}
	t.visitingPkgs[pkgName] = struct{}{}

	_, ok := t.allPkgs[pkgName]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, pkgName)
		return
	}

	tree, ok := t.packages[pkgName]
	if !ok {
		// package exists but not modified, so no rebuild needed
		fmt.Println("<= [skip] ", pkgName)

		delete(t.visitingPkgs, pkgName)
		t.builtPkgs[pkgName] = struct{}{}
		return
	}

	directUserImports := t.extractUserImports(tree)
	for _, imp := range directUserImports {
		t.buildPackage(imp)
	}

	ffiImports := t.extractFFIimports(tree)
	for _, imp := range ffiImports {
		t.buildFFIPackage(imp)
	}

	stdlibImports := t.extractStdLibImports(tree)
	for _, imp := range stdlibImports {
		t.buildStdLib(imp)
	}

	// Create new LLVM context for this package (Safe, as children are finished)
	llvm := NewLLVM(pkgName, t.outputDir)
	t.llvms[pkgName] = llvm

	// Resolve Imports: Declare symbols from direct and transitive dependencies (B and C)
	t.resolveUserImports(tree, directUserImports, llvm)
	t.resolveStdLibImports(tree, stdlibImports, llvm)
	llvm.AddImportEntry(state.PackageEntry{Name: pkgName, Alias: pkgName})

	// Compile
	fmt.Println("=> [compile] ", pkgName)
	t.compile(tree, llvm)

	// Dump
	llvm.Dump(t.outputDir, pkgName)

	// Mark as finished
	delete(t.visitingPkgs, pkgName)
	t.builtPkgs[pkgName] = struct{}{}
}

func (t *generator) buildFFIPackage(pkg state.PackageEntry) {
	if _, ok := t.ffiModules[pkg.Name]; ok {
		return
	}

	split := strings.Split(pkg.Name, ".")
	path := filepath.Join(t.outputDir, "tmp", fmt.Sprintf("%s.ll", split[len(split)-1]))

	// Read original .ll file
	input, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	// Parse LLVM IR
	mod, err := asm.ParseBytes(path, input)
	if err != nil {
		panic(err)
	}

	t.ffiModules[pkg.Name] = mod
}

func (t *generator) buildStdLib(pkg state.PackageEntry) {
	split := strings.Split(pkg.Name, ".")
	modName := split[len(split)-1]

	if _, ok := t.ffiModules[pkg.Name]; ok {
		return
	}

	path := filepath.Join(t.outputDir, "tmp", fmt.Sprintf("%s.ll", modName))

	// Read original .ll file
	input, err := os.ReadFile(path)
	if err != nil {
		_, ok := libs.ModuleList[modName]
		if ok {
			// indicates it has be overriden in codegen/bins
			return
		}

		panic(err)
	}

	// Parse LLVM IR
	mod, err := asm.ParseBytes(path, input)
	if err != nil {
		panic(err)
	}

	t.ffiModules[pkg.Name] = mod
}

// extractUserImports returns only the names of non-builtin imported packages.
func (t *generator) extractUserImports(tree ast.BlockStatement) []state.PackageEntry {
	var imports []state.PackageEntry
	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if !stc.IsBuiltIn() && !stc.IsFFI() {
				imports = append(imports, state.PackageEntry{Name: stc.Name, Alias: stc.Alias})
			}
		}
	}
	return imports
}

func (t *generator) extractFFIimports(tree ast.BlockStatement) []state.PackageEntry {
	var imports []state.PackageEntry
	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.IsFFI() {
				imports = append(imports, state.PackageEntry{Name: stc.Name, Alias: stc.Alias})
			}
		}
	}
	return imports
}

func (t *generator) extractStdLibImports(tree ast.BlockStatement) []state.PackageEntry {
	var imports []state.PackageEntry
	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.IsBuiltIn() {
				imports = append(imports, state.PackageEntry{Name: stc.Name, Alias: stc.Alias})
			}
		}
	}
	return imports
}

func (t *generator) resolveUserImports(tree ast.BlockStatement, directUserImports []state.PackageEntry, llvm *LLVM) {
	// declared map tracks all packages added to this module's symbol table to prevent redundancy
	// key should be alias name. resolveUserImports is called specific to a module & its imported
	// packages must be tracked with its alias names instead of fully qualified name.
	// e.g issue, imported package 'a' could have been imported multiple times in lower levels, tracking its
	// declaration by fully qualified name prevents running .Declare with alias. therefore track with
	// alias name
	declared := make(map[string]struct{})

	for _, st := range tree.Body {
		if stc, ok := st.(ast.ImportStatement); ok {
			if stc.IsBuiltIn() || stc.IsFFI() {
				modifiedFFIModuleName := stc.EndName()

				ffiModule, ok := t.ffiModules[stc.Name]
				if stc.IsBuiltIn() && !ok {
					// indicates an codegen/libs overriden module
					if _, ok := libs.ModuleList[modifiedFFIModuleName]; !ok {
						errorutils.Abort(errorutils.UnknownModule, stc.Name)
					}
				} else {
					RegisterDeclarations(llvm.st, state.PackageEntry{Name: modifiedFFIModuleName, Alias: stc.Alias}, llvm.GetModule(), ffiModule)
				}

				llvm.AddImportEntry(state.PackageEntry{Name: modifiedFFIModuleName, Alias: stc.Alias})
			} else {
				llvm.AddImportEntry(state.PackageEntry{Name: stc.Name, Alias: stc.Alias})
			}
		}
	}

	for _, pkgName := range directUserImports {
		t.recursiveTransitiveDeclaration(pkgName, llvm, declared)
	}
}

func (t *generator) resolveStdLibImports(tree ast.BlockStatement, stdLibImports []state.PackageEntry, llvm *LLVM) {
	for _, pkgName := range stdLibImports {
		splits := strings.Split(pkgName.Name, ".")
		t.importBasePackages(llvm.st.LibMethods, splits[len(splits)-1])
	}
}

// recursiveTransitiveDeclaration declares symbols of pkgName and all its dependencies (C)
// into the current module (A). This is the key to fixing the transitive dependency issue.
func (t *generator) recursiveTransitiveDeclaration(pkg state.PackageEntry, llvm *LLVM, declared map[string]struct{}) {
	pkgFullName := pkg.Name
	pkgAliasName := pkg.Alias

	if _, ok := declared[pkgAliasName]; ok {
		return
	}
	declared[pkgAliasName] = struct{}{}

	packageAST := ast.BlockStatement{}

	if _, ok := t.packages[pkgFullName]; !ok {
		// load from exports
		err := utils.LoadFromFile(filepath.Join(t.outputDir, fmt.Sprintf("%s.exports", pkgFullName)), &packageAST)
		if err != nil {
			panic(fmt.Errorf("error while loading exports: %v", err))
		}
	} else {
		packageAST = t.packages[pkgFullName]
	}

	subImports := t.extractUserImports(packageAST)
	for _, sub := range subImports {
		t.recursiveTransitiveDeclaration(sub, llvm, declared)
	}

	fmt.Printf("pkgFullName: %v\n", pkgFullName)
	for _, st := range packageAST.Body {
		if stc, ok := st.(ast.ImportStatement); ok && stc.IsBuiltIn() {
			t.importBasePackages(llvm.st.LibMethods, stc.EndName())
		}
	}

	// since i am tracking packages with alias names, this func might be called multiple times
	// for a package. .Declare() is assumed to avoid multiple llvm type/func declarations, otherwise
	// which is fatal.
	pipeline.NewPipeline(llvm.st, llvm.m, packageAST).Declare(pkg)
}

// importBasePackages resolve base module imports.
func (t *generator) importBasePackages(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		return
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s.%s", BUILTIN, module, name)
		methodMap[n] = f
	}
}

// compile passes tree through all steps to output IR
func (t *generator) compile(tree ast.BlockStatement, llvm *LLVM) {
	llvm.ParseAST(tree)
}

// LoadPackages loads package with their AST by going through project.ini file
// Packages will be named will modified by replacing '/' with '.' resulting
// in something like os.io instead of os/io.
func LoadPackages(projectDir string) (map[string]ast.BlockStatement, map[string]struct{}) {
	rootDir := projectDir

	modifiedPkgAST := make(map[string]ast.BlockStatement)
	allPkgImports := make(map[string]ast.BlockStatement)

	err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".pic" {
			return nil
		}

		tree := parser.ParseImports(path)

		// relative path from root
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// remove extension
		rel = strings.TrimSuffix(rel, ".pic")

		// normalize to forward slashes for package name
		pkgName := strings.ReplaceAll(filepath.ToSlash(rel), "/", ".")

		allPkgImports[pkgName] = tree
		return nil
	})

	if err != nil {
		panic(err)
	}

	cachedBuilder := tools.NewBuildCache(allPkgImports, projectDir)
	modifiedPkgs, err := cachedBuilder.CheckBuildCache()

	allPkgs := make(map[string]struct{})

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
		// relative path from root
		rel, err := filepath.Rel(rootDir, path)
		if err != nil {
			return err
		}

		// remove extension
		rel = strings.TrimSuffix(rel, ".pic")

		// normalize to forward slashes for package name
		pkgName := strings.ReplaceAll(filepath.ToSlash(rel), string(os.PathSeparator), ".")

		allPkgs[pkgName] = struct{}{}

		// ignore unmodified packages
		if _, ok := modifiedPkgs[pkgName]; !ok {
			return nil
		}
		tree := parser.ParseAll(path)

		modifiedPkgAST[pkgName] = tree
		return nil
	})

	if err != nil {
		panic(err)
	}

	utils.SaveToFile(filepath.Join(projectDir, BUILD, "build.meta"), tools.LastBuildTime{Time: time.Now()})
	return modifiedPkgAST, allPkgs
}
