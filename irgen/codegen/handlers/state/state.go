package state

import (
	"encoding/gob"
	"fmt"
	"os"
	"path"

	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/scope"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// State holds the global generator/interpreter state during IR generation.
type State struct {
	// OutputDir is where Info & IR files will be saved
	OutputDir string

	// Maintain a global type list to avoid redeclaring a llvm type
	GlobalTypeList map[string]types.Type
	// Maintain a global func list to avoid redeclaring a llvm func
	GlobalFuncList map[string]*ir.Func

	// current module name. fully qualified name.
	ModuleName string

	// Entry point function (usually "main")
	MainFunc *ir.Func

	// LLVM IR module
	Module *ir.Module

	// Handles type conversions and casts
	TypeHandler *tf.TypeHandler

	// Utility to create consistent identifiers
	// @todo: this is not properly used everywhere, doesn't support
	// all the requirement.
	// used to build identifier name which starts with current module fq name
	IdentifierBuilder *identifier.IdentifierBuilder

	// All global and scoped variables
	Vars *scope.VarTree

	// User-defined classes with metadata
	Classes map[string]*tf.MetaClass

	// Imported base library functions. comes from builtin module import.
	LibMethods map[string]function.Func

	// Global string counter for unique string literals
	// @todo: need to fix this
	StrCounter int

	// Class inheritance hierarchy
	TypeHeirarchy TypeHeirarchy

	// Garbage collector instance
	CI *c.Interface

	// loop
	Loopend []LoopEntry

	// imports
	Imports map[string]PackageEntry
}

func (t *State) LoadInfoMap(moduleName string) (map[string]ast.Type, error) {
	path := path.Join(t.OutputDir, fmt.Sprintf("%s.bin", moduleName))
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var infoMap map[string]ast.Type
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&infoMap)
	return infoMap, err
}

func (t *State) DumpInfo(outputDir string) {
	// currently I save return type of all functions
	// which is need to reconstruct ast.Type during foreign class inheritance.
	// @todo: in genenral I can maintain metadat struct that will be stored at the
	// end in info file.
	infoMap := make(map[string]ast.Type)
	for _, cls := range t.TypeHeirarchy.ClassDefs {
		for _, st := range cls.Body {
			if fn, ok := st.(ast.FunctionDefinitionStatement); ok {
				fqFnName := t.IdentifierBuilder.Attach(cls.Name, fn.Name)
				infoMap[fqFnName] = fn.ReturnType
			}
		}
	}

	binPath := fmt.Sprintf("%s/%s.bin", outputDir, t.ModuleName)
	if err := t.saveInfoMap(infoMap, binPath); err != nil {
		fmt.Printf("Error saving bin file: %v\n", err)
	}
}

func (t *State) saveInfoMap(infoMap map[string]ast.Type, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	return encoder.Encode(infoMap)
}

type PackageEntry struct {
	// represents fully qualified name of imported package. e.g, "os.io"
	Name string
	// alias of imported package
	Alias string
}

// LoopEntry is to keep track of loop end blocks, vital for the implementation
// of break & continue.
type LoopEntry struct {
	End *bc.BlockHolder
}

// TypeHeirarchy stores inheritance relationships between classes.
type TypeHeirarchy struct {
	// Parent class name for given class.
	Parent map[string]string
	// Child class ast for given class.
	Childs map[string][]ast.ClassDeclarationStatement

	// classes which doesn't inherit any other classes.
	Roots []ast.ClassDeclarationStatement

	// ast for classes in own module.
	ClassDefs map[string]ast.ClassDeclarationStatement

	// non root classes whose parent is foreign.
	Orphans []ast.ClassDeclarationStatement
}
