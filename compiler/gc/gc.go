package gc

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type GC struct {
	langAlloc   *ir.Func
	runtimeInit *ir.Func
}

var Instance *GC

func GetGC(mod *ir.Module) *GC {
	if Instance == nil {
		Instance = &GC{
			mod.NewFunc("lang_alloc", types.I8Ptr, ir.NewParam("", types.I64)),
			mod.NewFunc("runtime_init", types.Void),
		}

	}
	return Instance
}

func (t *GC) Init() *ir.Func {
	return t.runtimeInit
}

func (t *GC) Alloc() *ir.Func {
	return t.langAlloc
}
