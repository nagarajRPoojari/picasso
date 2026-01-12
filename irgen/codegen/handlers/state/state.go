package state

import (
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
	ClassRoots     []ast.ClassDeclarationStatement
	InterfaceRoots []ast.InterfaceDeclarationStatement
}

func NewTypeHeirarchy() *TypeHeirarchy {
	return &TypeHeirarchy{
		ClassRoots:     make([]ast.ClassDeclarationStatement, 0),
		InterfaceRoots: make([]ast.InterfaceDeclarationStatement, 0),
	}
}

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

	// List of interfaces
	Interfaces map[string]*tf.MetaInterface

	// Imported base library functions. comes from builtin module import.
	LibMethods map[string]function.Func

	// ffi method
	FFIModules map[string]FFIDeclarations

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

type FFIDeclarations struct {
	Methods map[string]*ir.Func
	Globals map[string]*ir.Global
	Types   map[string]*types.Type
}

type Handlers struct {
}

func NewCompileState(outputDir string, pkgName string, module *ir.Module) *State {
	return &State{
		OutputDir:         outputDir,
		GlobalTypeList:    make(map[string]types.Type),
		GlobalFuncList:    make(map[string]*ir.Func),
		ModuleName:        pkgName,
		Module:            module,
		TypeHandler:       tf.NewTypeHandler(),
		TypeHeirarchy:     *NewTypeHeirarchy(),
		FFIModules:        map[string]FFIDeclarations{},
		Vars:              scope.NewVarTree(),
		Classes:           make(map[string]*tf.MetaClass),
		Interfaces:        make(map[string]*tf.MetaInterface),
		IdentifierBuilder: identifier.NewIdentifierBuilder(pkgName),
		LibMethods:        make(map[string]function.Func),
		CI:                c.Instance,
		Imports:           make(map[string]PackageEntry),
	}
}
