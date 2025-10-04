package gc

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

const (
	ALLOC        = "lang_alloc"
	RUNTIME_INIT = "runtime_init"
	ARRAY_ALLOC  = "lang_alloc_array"
)

type GC struct {
	langAlloc   *ir.Func
	runtimeInit *ir.Func
	arrayAlloc  *ir.Func
}

var Instance *GC

func GetGC(mod *ir.Module) *GC {
	// must ensure the module doesn't already contains declaration,
	// to avoid redeclaring same functions.
	for _, f := range mod.Funcs {
		if f.Name() == ALLOC || f.Name() == RUNTIME_INIT || f.Name() == ARRAY_ALLOC {
			return Instance
		}
	}
	Instance = &GC{
		mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64)),
		mod.NewFunc(RUNTIME_INIT, types.Void),
		mod.NewFunc(ARRAY_ALLOC, types.NewPointer(types.NewStruct(types.I64, types.NewPointer(types.I8))), ir.NewParam("", types.I64), ir.NewParam("", types.I64)),
	}
	return Instance
}

func (t *GC) Init() *ir.Func {
	return t.runtimeInit
}

func (t *GC) Alloc() *ir.Func {
	return t.langAlloc
}

func (t *GC) ArrayAlloc() *ir.Func {
	return t.arrayAlloc
}
