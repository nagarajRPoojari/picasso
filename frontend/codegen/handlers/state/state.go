package state

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/frontend/ast"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/c"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/scope"
	function "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/func"
	tf "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

// State holds the global generator/interpreter state during IR generation.
type State struct {
	// Entry point function (usually "main")
	MainFunc *ir.Func
	// LLVM IR module
	Module *ir.Module

	// Handles type conversions and casts
	TypeHandler *tf.TypeHandler
	// Utility to create consistent identifiers
	IdentifierBuilder *identifier.IdentifierBuilder

	// All global and scoped variables
	Vars *scope.VarTree
	// User-defined classes with metadata
	Classes map[string]*tf.MetaClass
	// Imported library functions
	LibMethods map[string]function.Func

	// Global string counter for unique string literals
	StrCounter int
	// Class inheritance hierarchy
	TypeHeirarchy TypeHeirarchy
	// Garbage collector instance
	CI *c.Interface

	// loop
	Loopend []LoopEntry
}

type LoopEntry struct {
	End *bc.BlockHolder
}

// TypeHeirarchy stores inheritance relationships between classes.
type TypeHeirarchy struct {
	Parent    map[string]string
	Childs    map[string][]ast.ClassDeclarationStatement
	Roots     []ast.ClassDeclarationStatement
	ClassDefs map[string]ast.ClassDeclarationStatement
}
