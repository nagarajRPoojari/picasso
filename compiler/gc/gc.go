package gc

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

const (
	ALLOC        = "lang_alloc"
	RUNTIME_INIT = "runtime_init"
)

type GC struct {
	langAlloc   *ir.Func
	runtimeInit *ir.Func
}

var Instance *GC

func GetGC(mod *ir.Module) *GC {
	// must ensure the module doesn't already contains declaration,
	// to avoid redeclaring same functions.
	for _, f := range mod.Funcs {
		if f.Name() == ALLOC || f.Name() == RUNTIME_INIT {
			return Instance
		}
	}
	Instance = &GC{
		mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64)),
		mod.NewFunc(RUNTIME_INIT, types.Void),
	}
	return Instance
}

func (t *GC) Init() *ir.Func {
	return t.runtimeInit
}

func (t *GC) Alloc() *ir.Func {
	return t.langAlloc
}
