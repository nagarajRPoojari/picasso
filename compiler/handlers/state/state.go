package state

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/compiler/gc"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/identifier"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/scope"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type State struct {
	MainFunc *ir.Func
	// llvm module
	Module *ir.Module

	TypeHandler *tf.TypeHandler
	// utility tool to build consistent identifier names
	IdentifierBuilder *identifier.IdentifierBuilder

	// all global vars
	Vars *scope.VarTree
	// all methods including class methods & top level functions
	// custom classes defined by user
	Classes map[string]*tf.MetaClass

	// global string counter
	// @todo: move this to separate string module
	StrCounter int

	LibMethods map[string]function.Func

	GC *gc.GC
}
