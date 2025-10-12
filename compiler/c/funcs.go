package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (t *Interface) registerFuncs(mod *ir.Module) {
	// custom runtime
	t.Funcs[ALLOC] = mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	t.Funcs[RUNTIME_INIT] = mod.NewFunc(RUNTIME_INIT, types.Void)

	t.Funcs[ARRAY_ALLOC] = mod.NewFunc(ARRAY_ALLOC, types.NewPointer(Array), ir.NewParam("", types.I64), ir.NewParam("", types.I64))

	t.Funcs[PRINTF] = mod.NewFunc(PRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[PRINTF].Sig.Variadic = true

	t.Funcs[MALLOC] = mod.NewFunc(MALLOC, types.I8Ptr, ir.NewParam("size", types.I64))

	t.Funcs[FREE] = mod.NewFunc(FREE, types.Void, ir.NewParam("ptr", types.I8Ptr))

	t.Funcs[STRLEN] = mod.NewFunc(STRLEN, types.I64, ir.NewParam("s", types.I8Ptr))

	t.Funcs[MEMCPY] = mod.NewFunc(MEMCPY, types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr),
		ir.NewParam("n", types.I64),
	)

	t.Funcs[MEMSET] = mod.NewFunc(MEMSET, types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("val", types.I32),
		ir.NewParam("n", types.I64),
	)

	t.Funcs[MEMMOVE] = mod.NewFunc(MEMMOVE, types.I8Ptr,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr),
		ir.NewParam("n", types.I64),
	)

	t.Funcs[EXIT] = mod.NewFunc(EXIT, types.Void, ir.NewParam("code", types.I32))
}
