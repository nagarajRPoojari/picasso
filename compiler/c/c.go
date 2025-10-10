package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
)

const (
	ALLOC        = "lang_alloc"
	RUNTIME_INIT = "runtime_init"
	ARRAY_ALLOC  = "lang_alloc_array"
)

type Interface struct {
	langAlloc   *ir.Func
	runtimeInit *ir.Func
	arrayAlloc  *ir.Func
}

var Instance *Interface

var ARRAYSTRUCT = types.NewStruct(
	types.I64,                   // length
	types.NewPointer(types.I8),  // data
	types.NewPointer(types.I64), // shape (i64*)
	types.I64,                   // rank
)

func init() {
	ARRAYSTRUCT.SetName(constants.ARRAY)
}

func NewInterface(mod *ir.Module) *Interface {
	// must ensure the module doesn't already contains declaration,
	// to avoid redeclaring same functions.
	for _, f := range mod.Funcs {
		if f.Name() == ALLOC || f.Name() == RUNTIME_INIT || f.Name() == ARRAY_ALLOC {
			return Instance
		}
	}
	Instance = &Interface{
		mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64)),
		mod.NewFunc(RUNTIME_INIT, types.Void),
		mod.NewFunc(ARRAY_ALLOC, types.NewPointer(ARRAYSTRUCT), ir.NewParam("", types.I64), ir.NewParam("", types.I64)),
	}
	return Instance
}

func (t *Interface) Init() *ir.Func {
	return t.runtimeInit
}

func (t *Interface) Alloc() *ir.Func {
	return t.langAlloc
}

func (t *Interface) ArrayAlloc() *ir.Func {
	return t.arrayAlloc
}
