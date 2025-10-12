package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func (t *Interface) registerFuncs(mod *ir.Module) {
	// custom runtime
	t.initRuntime(mod)
	// stdio
	t.initStdio(mod)
}

func (t *Interface) initRuntime(mod *ir.Module) {
	// @alloc
	t.Funcs[ALLOC] = mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	// @runtime_init
	t.Funcs[RUNTIME_INIT] = mod.NewFunc(RUNTIME_INIT, types.Void)

	// @array_alloc
	t.Funcs[ARRAY_ALLOC] = mod.NewFunc(ARRAY_ALLOC, types.NewPointer(Array), ir.NewParam("", types.I64), ir.NewParam("", types.I64))
}

func (t *Interface) initStdio(mod *ir.Module) {
	// @printf
	t.Funcs[PRINTF] = mod.NewFunc(PRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[PRINTF].Sig.Variadic = true

	// @scanf
	t.Funcs[SCANF] = mod.NewFunc(SCANF, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[SCANF].Sig.Variadic = true

	// @fopen
	t.Funcs[FOPEN] = mod.NewFunc(FOPEN, types.I8Ptr,
		ir.NewParam("filename", types.I8Ptr),
		ir.NewParam("mode", types.I8Ptr),
	)

	// @fclose
	t.Funcs[FCLOSE] = mod.NewFunc(FCLOSE, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fprintf
	t.Funcs[FPRINTF] = mod.NewFunc(FPRINTF, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("format", types.I8Ptr),
	)
	t.Funcs[FPRINTF].Sig.Variadic = true

	// @fscanf
	t.Funcs[FSCANF] = mod.NewFunc(FSCANF, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("format", types.I8Ptr),
	)
	t.Funcs[FSCANF].Sig.Variadic = true

	// @fputs
	t.Funcs[FPUTS] = mod.NewFunc(FPUTS, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("file", types.I8Ptr),
	)

	// @fgets
	t.Funcs[FGETS] = mod.NewFunc(FGETS, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("file", types.I8Ptr),
	)

	// @fflush
	t.Funcs[FFLUSH] = mod.NewFunc(FFLUSH, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fseek
	t.Funcs[FSEEK] = mod.NewFunc(FSEEK, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("offset", types.I64),
		ir.NewParam("whence", types.I32),
	)
}
