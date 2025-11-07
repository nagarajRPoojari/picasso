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
	// atomic
	t.initAtomicFuncs(mod)
}

func (t *Interface) initRuntime(mod *ir.Module) {

	t.Funcs[HASH] = mod.NewFunc(HASH, types.I64, ir.NewParam("data", types.I8Ptr), ir.NewParam("len", types.I64))

	t.Funcs[STRLEN] = mod.NewFunc(STRLEN, types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	t.Funcs[STRCMP] = mod.NewFunc(STRCMP, types.I32, ir.NewParam("", types.NewPointer(types.I8)), ir.NewParam("", types.NewPointer(types.I8)))

	t.Funcs[MEMCPY] = mod.NewFunc("llvm.memcpy.p0i8.p0i8.i64",
		types.Void,
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("src", types.I8Ptr),
		ir.NewParam("len", types.I64),
		ir.NewParam("isvolatile", types.I1),
	)
	// @thread
	fnType := types.NewFunc(
		types.NewPointer(types.I8),
		types.NewPointer(types.I8),
	)
	t.Funcs[THREAD] = mod.NewFunc(THREAD, types.Void,
		ir.NewParam("", types.NewPointer(fnType)),
	)

	// @self_yield
	t.Funcs[SELF_YIELD] = mod.NewFunc(SELF_YIELD, types.Void)

	// @alloc
	t.Funcs[ALLOC] = mod.NewFunc(ALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	// @runtime_init
	t.Funcs[RUNTIME_INIT] = mod.NewFunc(RUNTIME_INIT, types.Void)

	// @array_alloc
	t.Funcs[ARRAY_ALLOC] = mod.NewFunc(ARRAY_ALLOC, types.NewPointer(t.Types[ARRAY]), ir.NewParam("", types.I64), ir.NewParam("", types.I64))
}

func (t *Interface) initAtomicFuncs(mod *ir.Module) {
	// --- @bool ---
	t.Funcs[ATOMIC_STORE_BOOL] = mod.NewFunc(ATOMIC_STORE_BOOL,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_BOOL])),
		ir.NewParam("val", types.I1),
	)
	t.Funcs[ATOMIC_LOAD_BOOL] = mod.NewFunc(ATOMIC_LOAD_BOOL,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_BOOL])),
	)

	// --- @int8 ---
	t.Funcs[ATOMIC_STORE_CHAR] = mod.NewFunc(ATOMIC_STORE_CHAR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_STORE_INT8] = t.Funcs[ATOMIC_STORE_CHAR]
	t.Funcs[ATOMIC_LOAD_CHAR] = mod.NewFunc(ATOMIC_LOAD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_CHAR])),
	)
	t.Funcs[ATOMIC_LOAD_INT8] = t.Funcs[ATOMIC_LOAD_CHAR]
	t.Funcs[ATOMIC_ADD_CHAR] = mod.NewFunc(ATOMIC_ADD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_ADD_INT8] = t.Funcs[ATOMIC_ADD_CHAR]
	t.Funcs[ATOMIC_SUB_CHAR] = mod.NewFunc(ATOMIC_SUB_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_SUB_INT8] = t.Funcs[ATOMIC_SUB_CHAR]

	// --- @int16 ---
	t.Funcs[ATOMIC_STORE_SHORT] = mod.NewFunc(ATOMIC_STORE_SHORT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_STORE_INT16] = t.Funcs[ATOMIC_STORE_SHORT]
	t.Funcs[ATOMIC_LOAD_SHORT] = mod.NewFunc(ATOMIC_LOAD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_SHORT])),
	)
	t.Funcs[ATOMIC_LOAD_INT16] = t.Funcs[ATOMIC_LOAD_SHORT]
	t.Funcs[ATOMIC_ADD_SHORT] = mod.NewFunc(ATOMIC_ADD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_ADD_INT16] = t.Funcs[ATOMIC_ADD_SHORT]
	t.Funcs[ATOMIC_SUB_SHORT] = mod.NewFunc(ATOMIC_SUB_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_SUB_INT16] = t.Funcs[ATOMIC_SUB_SHORT]

	// --- @int32 ---
	t.Funcs[ATOMIC_STORE_INT] = mod.NewFunc(ATOMIC_STORE_INT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_STORE_INT32] = t.Funcs[ATOMIC_STORE_INT]
	t.Funcs[ATOMIC_LOAD_INT] = mod.NewFunc(ATOMIC_LOAD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_INT])),
	)
	t.Funcs[ATOMIC_LOAD_INT32] = t.Funcs[ATOMIC_LOAD_INT]
	t.Funcs[ATOMIC_ADD_INT] = mod.NewFunc(ATOMIC_ADD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_ADD_INT32] = t.Funcs[ATOMIC_ADD_INT]
	t.Funcs[ATOMIC_SUB_INT] = mod.NewFunc(ATOMIC_SUB_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_SUB_INT32] = t.Funcs[ATOMIC_SUB_INT]

	// --- @int64 ---
	t.Funcs[ATOMIC_STORE_LONG] = mod.NewFunc(ATOMIC_STORE_LONG,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_STORE_INT64] = t.Funcs[ATOMIC_STORE_LONG]
	t.Funcs[ATOMIC_LOAD_LONG] = mod.NewFunc(ATOMIC_LOAD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_LONG])),
	)
	t.Funcs[ATOMIC_LOAD_INT64] = t.Funcs[ATOMIC_LOAD_LONG]
	t.Funcs[ATOMIC_ADD_LONG] = mod.NewFunc(ATOMIC_ADD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_ADD_INT64] = t.Funcs[ATOMIC_ADD_LONG]
	t.Funcs[ATOMIC_SUB_LONG] = mod.NewFunc(ATOMIC_SUB_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_SUB_INT64] = t.Funcs[ATOMIC_SUB_LONG]

	// --- floats and others unchanged ---
	t.Funcs[ATOMIC_STORE_FLOAT] = mod.NewFunc(ATOMIC_STORE_FLOAT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_FLOAT])),
		ir.NewParam("val", types.Float),
	)
	t.Funcs[ATOMIC_LOAD_FLOAT] = mod.NewFunc(ATOMIC_LOAD_FLOAT,
		types.Float,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_FLOAT])),
	)

	t.Funcs[ATOMIC_STORE_DOUBLE] = mod.NewFunc(ATOMIC_STORE_DOUBLE,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_DOUBLE])),
		ir.NewParam("val", types.Double),
	)
	t.Funcs[ATOMIC_LOAD_DOUBLE] = mod.NewFunc(ATOMIC_LOAD_DOUBLE,
		types.Double,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_DOUBLE])),
	)

	t.Funcs[ATOMIC_STORE_PTR] = mod.NewFunc(ATOMIC_STORE_PTR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_PTR])),
		ir.NewParam("val", types.NewPointer(types.I8)),
	)
	t.Funcs[ATOMIC_LOAD_PTR] = mod.NewFunc(ATOMIC_LOAD_PTR,
		types.NewPointer(types.I8),
		ir.NewParam("ptr", types.NewPointer(t.Types[ATOMIC_PTR])),
	)
}

func (t *Interface) initStdio(mod *ir.Module) {
	// @printf
	t.Funcs[PRINTF] = mod.NewFunc(PRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[PRINTF].Sig.Variadic = true

	t.Funcs[APRINTF] = mod.NewFunc(APRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[APRINTF].Sig.Variadic = true

	// @scanf
	t.Funcs[SCANF] = mod.NewFunc(SCANF, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[SCANF].Sig.Variadic = true

	t.Funcs[ASCAN] = mod.NewFunc(ASCAN, types.I8Ptr, ir.NewParam("size", types.I64))

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

	// @afread
	t.Funcs[AFREAD] = mod.NewFunc(AFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// @afwrite
	t.Funcs[AFWRITE] = mod.NewFunc(AFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
}
