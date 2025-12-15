package generator

import (
	"fmt"
	"os"

	"github.com/go-ini/ini"
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	errorutils "github.com/nagarajRPoojari/niyama/frontend/codegen/error"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/constants"
	"github.com/nagarajRPoojari/niyama/frontend/parser"
)

type generator struct {
	llvm     *LLVM
	packages map[string]ast.BlockStatement
	tree     ast.BlockStatement
}

func NewGenerator() *generator {
	return &generator{
		llvm: NewLLVM(),
	}
}

func (t *generator) Build(path string) {
	m := LoadPackages(path)
	t.packages = m
	t.SetAST(t.BuildUniModule())
}

func (t *generator) SetAST(tree ast.BlockStatement) {
	t.tree = tree
}

func (t *generator) Compile() {
	t.llvm.ParseAST(t.tree)
}

func (t *generator) BuildUniModule() ast.BlockStatement {
	mainModule := t.packages[constants.MAIN]
	imported := make(map[string]struct{})
	stack := make(map[string]struct{})
	t.resolveImportsRecursive(&mainModule, imported, stack)
	return mainModule
}

// resolveImportsRecursive traverses imports in `module` and inlines their bodies.
func (t *generator) resolveImportsRecursive(module *ast.BlockStatement, imported map[string]struct{}, stack map[string]struct{}) {
	for i := 0; i < len(module.Body); i++ {
		st := module.Body[i]

		importStmt, ok := st.(ast.ImportStatement)
		if !ok {
			continue
		}

		if importStmt.From != importStmt.Name {
			if importStmt.From == constants.BUILTIN {
				continue
			}
			errorutils.Abort(errorutils.InvalidModulerSource, importStmt.From, importStmt.Name)
		}

		pkgName := importStmt.Name
		if _, seen := stack[pkgName]; seen {
			panic("recursive import detected: " + pkgName)
		}

		if _, seen := imported[pkgName]; seen {
			// already imported
			continue
		}

		pkgModule, exists := t.packages[pkgName]
		if !exists {
			// might be a native library
			continue
		}

		stack[pkgName] = struct{}{}
		imported[pkgName] = struct{}{}
		t.resolveImportsRecursive(&pkgModule, imported, stack)

		module.Body = append(
			module.Body[:i],
			append(pkgModule.Body, module.Body[i+1:]...)...,
		)

		i += len(pkgModule.Body) - 1

		delete(stack, pkgName)
	}
}

func (t *generator) Dump(file string) {
	t.llvm.Dump(file)
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
