package tools

import (
	"fmt"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
)

type ImportTree struct {
	graph map[string][]string
}

type BuildCache struct {
	pkgs           map[string]ast.BlockStatement
	projectRootDir string
}

func NewBuildCache(pkgs map[string]ast.BlockStatement, projectRootDir string) *BuildCache {
	buildCache := &BuildCache{
		pkgs:           pkgs,
		projectRootDir: projectRootDir,
	}
	return buildCache
}

func (t *BuildCache) buildTree(it *ImportTree, visitStack map[string]struct{}, pkgName string) error {
	if _, ok := visitStack[pkgName]; ok {
		return fmt.Errorf("cyclic dependency found involving %s", pkgName)
	}
	pkg := t.pkgs[pkgName]
	it.graph[pkgName] = make([]string, 0)
	visitStack[pkgName] = struct{}{}
	for _, stmt := range pkg.Body {
		if imp, ok := stmt.(ast.ImportStatement); ok {
			// avoid base package
			if imp.IsBasePkg() {
				continue
			}
			// use fully qualified name .Name instead of .Alias
			it.graph[pkgName] = append(it.graph[pkgName], imp.Name)
			t.buildTree(it, visitStack, imp.Name)
		}
	}
	return nil
}

// func (t *BuildCache)
