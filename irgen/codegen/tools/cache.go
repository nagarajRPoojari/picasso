package tools

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/utils"
)

type LastBuildTime struct {
	Time time.Time
}

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
			if imp.IsBuiltIn() || imp.IsFFI() {
				continue
			}
			// use fully qualified name .Name instead of .Alias
			it.graph[pkgName] = append(it.graph[pkgName], imp.Name)
			t.buildTree(it, visitStack, imp.Name)
		}
	}
	return nil
}

func (t *BuildCache) pkgNameFromPath(root, path string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return ""
	}

	rel = strings.TrimSuffix(rel, filepath.Ext(rel))
	parts := strings.Split(rel, string(os.PathSeparator))
	return strings.Join(parts, ".")
}

func (t *BuildCache) findModifiedPkgs(lastBuild time.Time) (map[string]struct{}, error) {
	modified := make(map[string]struct{})

	err := filepath.WalkDir(t.projectRootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip build directory
		if d.IsDir() && d.Name() == "build" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		if filepath.Ext(d.Name()) != ".pic" {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		if info.ModTime().After(lastBuild) {
			pkg := t.pkgNameFromPath(t.projectRootDir, path)
			modified[pkg] = struct{}{}
		}

		return nil
	})

	return modified, err
}

func (t *BuildCache) reverseImportTree(it *ImportTree) map[string][]string {
	rev := make(map[string][]string)

	for pkg, imports := range it.graph {
		for _, imp := range imports {
			rev[imp] = append(rev[imp], pkg)
		}
	}

	return rev
}

func (t *BuildCache) propagateModified(initial map[string]struct{}, reverse map[string][]string) map[string]struct{} {

	dirty := make(map[string]struct{})
	queue := make([]string, 0)

	for pkg := range initial {
		dirty[pkg] = struct{}{}
		queue = append(queue, pkg)
	}

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]

		for _, parent := range reverse[p] {
			if _, seen := dirty[parent]; seen {
				continue
			}
			dirty[parent] = struct{}{}
			queue = append(queue, parent)
		}
	}

	return dirty
}

func (t *BuildCache) CheckBuildCache() (map[string]struct{}, error) {
	lastBuild := LastBuildTime{}
	err := utils.LoadFromFile(filepath.Join(t.projectRootDir, "build", "build.meta"), &lastBuild)
	if err != nil {
	}

	modified, err := t.findModifiedPkgs(lastBuild.Time)

	if err != nil {
		return nil, err
	}

	it := &ImportTree{
		graph: make(map[string][]string),
	}

	visit := make(map[string]struct{})
	if err := t.buildTree(it, visit, "start"); err != nil {
		return nil, err
	}

	reverse := t.reverseImportTree(it)
	dirty := t.propagateModified(modified, reverse)

	return dirty, nil
}
