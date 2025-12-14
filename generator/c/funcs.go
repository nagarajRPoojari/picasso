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

	t.initStrs(mod)
}

func (t *Interface) initRuntime(mod *ir.Module) {

	t.Funcs[FUNC_HASH] = mod.NewFunc(FUNC_HASH, types.I64, ir.NewParam("data", types.I8Ptr), ir.NewParam("len", types.I64))

	t.Funcs[FUNC_STRLEN] = mod.NewFunc(FUNC_STRLEN, types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	t.Funcs[FUNC_STRCMP] = mod.NewFunc(FUNC_STRCMP, types.I32, ir.NewParam("", types.NewPointer(types.I8)), ir.NewParam("", types.NewPointer(types.I8)))

	t.Funcs[FUNC_MEMCPY] = mod.NewFunc("llvm.memcpy.p0i8.p0i8.i64",
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
	t.Funcs[FUNC_THREAD] = mod.NewFunc(FUNC_THREAD, types.Void,
		ir.NewParam("", types.NewPointer(fnType)),
	)

	// @self_yield
	t.Funcs[FUNC_SELF_YIELD] = mod.NewFunc(FUNC_SELF_YIELD, types.Void)

	// @alloc
	t.Funcs[FUNC_ALLOC] = mod.NewFunc(FUNC_ALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	t.Funcs[FUNC_MALLOC] = mod.NewFunc(FUNC_MALLOC, types.I8Ptr, ir.NewParam("", types.I64))

	// @runtime_init
	t.Funcs[FUNC_RUNTIME_INIT] = mod.NewFunc(FUNC_RUNTIME_INIT, types.Void)

	// @array_alloc
	t.Funcs[FUNC_ARRAY_ALLOC] = mod.NewFunc(FUNC_ARRAY_ALLOC, types.NewPointer(t.Types[TYPE_ARRAY]), ir.NewParam("", types.I32), ir.NewParam("", types.I32), ir.NewParam("", types.I32))

	t.Funcs[__UTILS__FUNC_DEBUG_ARRAY_INFO] = mod.NewFunc(
		__UTILS__FUNC_DEBUG_ARRAY_INFO,
		types.Void, ir.NewParam("", types.NewPointer(t.Types[TYPE_ARRAY])),
	)
}

func (t *Interface) initAtomicFuncs(mod *ir.Module) {
	// --- @bool ---
	t.Funcs[ATOMIC_STORE_BOOL] = mod.NewFunc(ATOMIC_STORE_BOOL,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
		ir.NewParam("val", types.I1),
	)
	t.Funcs[ATOMIC_LOAD_BOOL] = mod.NewFunc(ATOMIC_LOAD_BOOL,
		types.I1,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_BOOL])),
	)

	// --- @int8 ---
	t.Funcs[ATOMIC_STORE_CHAR] = mod.NewFunc(ATOMIC_STORE_CHAR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_STORE_INT8] = t.Funcs[ATOMIC_STORE_CHAR]
	t.Funcs[ATOMIC_LOAD_CHAR] = mod.NewFunc(ATOMIC_LOAD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
	)
	t.Funcs[ATOMIC_LOAD_INT8] = t.Funcs[ATOMIC_LOAD_CHAR]
	t.Funcs[ATOMIC_ADD_CHAR] = mod.NewFunc(ATOMIC_ADD_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_ADD_INT8] = t.Funcs[ATOMIC_ADD_CHAR]
	t.Funcs[ATOMIC_SUB_CHAR] = mod.NewFunc(ATOMIC_SUB_CHAR,
		types.I8,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_CHAR])),
		ir.NewParam("val", types.I8),
	)
	t.Funcs[ATOMIC_SUB_INT8] = t.Funcs[ATOMIC_SUB_CHAR]

	// --- @int16 ---
	t.Funcs[ATOMIC_STORE_SHORT] = mod.NewFunc(ATOMIC_STORE_SHORT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_STORE_INT16] = t.Funcs[ATOMIC_STORE_SHORT]
	t.Funcs[ATOMIC_LOAD_SHORT] = mod.NewFunc(ATOMIC_LOAD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
	)
	t.Funcs[ATOMIC_LOAD_INT16] = t.Funcs[ATOMIC_LOAD_SHORT]
	t.Funcs[ATOMIC_ADD_SHORT] = mod.NewFunc(ATOMIC_ADD_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_ADD_INT16] = t.Funcs[ATOMIC_ADD_SHORT]
	t.Funcs[ATOMIC_SUB_SHORT] = mod.NewFunc(ATOMIC_SUB_SHORT,
		types.I16,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_SHORT])),
		ir.NewParam("val", types.I16),
	)
	t.Funcs[ATOMIC_SUB_INT16] = t.Funcs[ATOMIC_SUB_SHORT]

	// --- @int32 ---
	t.Funcs[ATOMIC_STORE_INT] = mod.NewFunc(ATOMIC_STORE_INT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_STORE_INT32] = t.Funcs[ATOMIC_STORE_INT]
	t.Funcs[ATOMIC_LOAD_INT] = mod.NewFunc(ATOMIC_LOAD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
	)
	t.Funcs[ATOMIC_LOAD_INT32] = t.Funcs[ATOMIC_LOAD_INT]
	t.Funcs[ATOMIC_ADD_INT] = mod.NewFunc(ATOMIC_ADD_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_ADD_INT32] = t.Funcs[ATOMIC_ADD_INT]
	t.Funcs[ATOMIC_SUB_INT] = mod.NewFunc(ATOMIC_SUB_INT,
		types.I32,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_INT])),
		ir.NewParam("val", types.I32),
	)
	t.Funcs[ATOMIC_SUB_INT32] = t.Funcs[ATOMIC_SUB_INT]

	// --- @int64 ---
	t.Funcs[ATOMIC_STORE_LONG] = mod.NewFunc(ATOMIC_STORE_LONG,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_STORE_INT64] = t.Funcs[ATOMIC_STORE_LONG]
	t.Funcs[ATOMIC_LOAD_LONG] = mod.NewFunc(ATOMIC_LOAD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
	)
	t.Funcs[ATOMIC_LOAD_INT64] = t.Funcs[ATOMIC_LOAD_LONG]
	t.Funcs[ATOMIC_ADD_LONG] = mod.NewFunc(ATOMIC_ADD_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_ADD_INT64] = t.Funcs[ATOMIC_ADD_LONG]
	t.Funcs[ATOMIC_SUB_LONG] = mod.NewFunc(ATOMIC_SUB_LONG,
		types.I64,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_LONG])),
		ir.NewParam("val", types.I64),
	)
	t.Funcs[ATOMIC_SUB_INT64] = t.Funcs[ATOMIC_SUB_LONG]

	// --- floats and others unchanged ---
	t.Funcs[ATOMIC_STORE_FLOAT] = mod.NewFunc(ATOMIC_STORE_FLOAT,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT])),
		ir.NewParam("val", types.Float),
	)
	t.Funcs[ATOMIC_LOAD_FLOAT] = mod.NewFunc(ATOMIC_LOAD_FLOAT,
		types.Float,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_FLOAT])),
	)

	t.Funcs[ATOMIC_STORE_DOUBLE] = mod.NewFunc(ATOMIC_STORE_DOUBLE,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_DOUBLE])),
		ir.NewParam("val", types.Double),
	)
	t.Funcs[ATOMIC_LOAD_DOUBLE] = mod.NewFunc(ATOMIC_LOAD_DOUBLE,
		types.Double,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_DOUBLE])),
	)

	t.Funcs[ATOMIC_STORE_PTR] = mod.NewFunc(ATOMIC_STORE_PTR,
		types.Void,
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
		ir.NewParam("val", types.NewPointer(types.I8)),
	)
	t.Funcs[ATOMIC_LOAD_PTR] = mod.NewFunc(ATOMIC_LOAD_PTR,
		types.NewPointer(types.I8),
		ir.NewParam("ptr", types.NewPointer(t.Types[TYPE_ATOMIC_PTR])),
	)
}

func (t *Interface) initStdio(mod *ir.Module) {
	// @fopen
	t.Funcs[FUNC_FOPEN] = mod.NewFunc(FUNC_FOPEN, types.I8Ptr,
		ir.NewParam("filename", types.I8Ptr),
		ir.NewParam("mode", types.I8Ptr),
	)
	// @fclose
	t.Funcs[FUNC_FCLOSE] = mod.NewFunc(FUNC_FCLOSE, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fflush
	t.Funcs[FUNC_FFLUSH] = mod.NewFunc(FUNC_FFLUSH, types.I32,
		ir.NewParam("stream", types.I8Ptr),
	)

	// @fseek
	t.Funcs[FUNC_FSEEK] = mod.NewFunc(FUNC_FSEEK, types.I32,
		ir.NewParam("stream", types.I8Ptr),
		ir.NewParam("offset", types.I64),
		ir.NewParam("whence", types.I32),
	)

	// @aprintf
	t.Funcs[FUNC_APRINTF] = mod.NewFunc(FUNC_APRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[FUNC_APRINTF].Sig.Variadic = true
	// @sprintf
	t.Funcs[FUNC_SPRINTF] = mod.NewFunc(FUNC_SPRINTF, types.I32, ir.NewParam("", types.I8Ptr))
	t.Funcs[FUNC_SPRINTF].Sig.Variadic = true

	// @ascanf
	t.Funcs[FUNC_ASCAN] = mod.NewFunc(FUNC_ASCAN, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_ASCAN].Sig.Variadic = true
	// @sscanf
	t.Funcs[FUNC_SSCAN] = mod.NewFunc(FUNC_SSCAN, types.I32, ir.NewParam("format", types.I8Ptr))
	t.Funcs[FUNC_SSCAN].Sig.Variadic = true

	// @afread
	t.Funcs[FUNC_AFREAD] = mod.NewFunc(FUNC_AFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfreed
	t.Funcs[FUNC_SFREAD] = mod.NewFunc(FUNC_SFREAD, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)

	// @afwrite
	t.Funcs[FUNC_AFWRITE] = mod.NewFunc(FUNC_AFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
	// @sfwrite
	t.Funcs[FUNC_SFWRITE] = mod.NewFunc(FUNC_SFWRITE, types.I32,
		ir.NewParam("fd", types.I8Ptr),
		ir.NewParam("dest", types.I8Ptr),
		ir.NewParam("n", types.I64),
		ir.NewParam("offset", types.I64),
	)
}

func (t *Interface) initStrs(mod *ir.Module) {
	// @format
	t.Funcs[FUNC_FORMAT] = mod.NewFunc(FUNC_FORMAT, types.I8Ptr,
		ir.NewParam("fmt", types.I8Ptr),
	)

	// @len
	t.Funcs[FUNC_LEN] = mod.NewFunc(FUNC_LEN, types.I32,
		ir.NewParam("str", types.I8Ptr),
	)

	// @compare
	t.Funcs[FUNC_COMPARE] = mod.NewFunc(FUNC_COMPARE, types.I32,
		ir.NewParam("a", types.I8Ptr),
		ir.NewParam("b", types.I8Ptr),
	)

}
