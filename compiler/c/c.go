package c

import (
	"sync"

	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
)

const (
	THREAD     = "thread"
	SELF_YIELD = "self_yield"

	// async I/O calls

	// blocking I/O calls
	PRINTF  = "printf"
	SCANF   = "scanf"
	FPRINTF = "fprintf"
	FSCANF  = "fscanf"
	FPUTS   = "fputs"
	FGETS   = "fgets"
	FOPEN   = "fopen"
	FCLOSE  = "fclose"
	FFLUSH  = "fflush"
	FSEEK   = "fseek"

	// memory calls
	MALLOC       = "malloc"
	FREE         = "free"
	MEMCPY       = "memcpy"
	MEMSET       = "memset"
	MEMMOVE      = "memmove"
	ALLOC        = "lang_alloc"
	RUNTIME_INIT = "runtime_init"
	ARRAY_ALLOC  = "lang_alloc_array"
	ARRAY        = "array"

	// string
	STRLEN = "strlen"

	EXIT = "exit"

	// sync calls @todo
	// atomic base operations
	ATOMIC_STORE = "atomic_store"
	ATOMIC_LOAD  = "atomic_load"
	ATOMIC_ADD   = "atomic_add"
	ATOMIC_SUB   = "atomic_sub"

	ATOMIC_STORE_BOOL = "atomic_store_bool"
	ATOMIC_LOAD_BOOL  = "atomic_load_bool"

	ATOMIC_STORE_CHAR = "atomic_store_char"
	ATOMIC_STORE_INT8 = "atomic_store_int8" // @alias
	ATOMIC_LOAD_CHAR  = "atomic_load_char"
	ATOMIC_LOAD_INT8  = "atomic_load_int8" // @alias
	ATOMIC_ADD_CHAR   = "atomic_add_char"
	ATOMIC_ADD_INT8   = "atomic_add_int8" // @alias
	ATOMIC_SUB_CHAR   = "atomic_sub_char"
	ATOMIC_SUB_INT8   = "atomic_sub_int8" // @alias

	ATOMIC_STORE_SHORT = "atomic_store_short"
	ATOMIC_STORE_INT16 = "atomic_store_int16" // @alias
	ATOMIC_LOAD_SHORT  = "atomic_load_short"
	ATOMIC_LOAD_INT16  = "atomic_load_int16" // @alias
	ATOMIC_ADD_SHORT   = "atomic_add_short"
	ATOMIC_ADD_INT16   = "atomic_add_int16" // @alias
	ATOMIC_SUB_SHORT   = "atomic_sub_short"
	ATOMIC_SUB_INT16   = "atomic_sub_int16" // @alias

	ATOMIC_STORE_INT   = "atomic_store_int"
	ATOMIC_STORE_INT32 = "atomic_store_int32" // @alias
	ATOMIC_LOAD_INT    = "atomic_load_int"
	ATOMIC_LOAD_INT32  = "atomic_load_int32" // @alias
	ATOMIC_ADD_INT     = "atomic_add_int"
	ATOMIC_ADD_INT32   = "atomic_add_int32" // @alias
	ATOMIC_SUB_INT     = "atomic_sub_int"
	ATOMIC_SUB_INT32   = "atomic_sub_int32" // @alias

	ATOMIC_STORE_LONG  = "atomic_store_long"
	ATOMIC_STORE_INT64 = "atomic_store_int64" // @alias
	ATOMIC_LOAD_LONG   = "atomic_load_long"
	ATOMIC_LOAD_INT64  = "atomic_load_int64" // @alias
	ATOMIC_ADD_LONG    = "atomic_add_long"
	ATOMIC_ADD_INT64   = "atomic_add_int64" // @alias
	ATOMIC_SUB_LONG    = "atomic_sub_long"
	ATOMIC_SUB_INT64   = "atomic_sub_int64" // @alias

	ATOMIC_STORE_LLONG = "atomic_store_llong"
	ATOMIC_LOAD_LLONG  = "atomic_load_llong"
	ATOMIC_ADD_LLONG   = "atomic_add_llong"
	ATOMIC_SUB_LLONG   = "atomic_sub_llong"

	ATOMIC_STORE_FLOAT = "atomic_store_float"
	ATOMIC_LOAD_FLOAT  = "atomic_load_float"

	ATOMIC_STORE_DOUBLE = "atomic_store_double"
	ATOMIC_LOAD_DOUBLE  = "atomic_load_double"

	ATOMIC_STORE_PTR = "atomic_store_ptr"
	ATOMIC_LOAD_PTR  = "atomic_load_ptr"

	ATOMIC_BOOL = "atomic_bool_t"

	ATOMIC_CHAR = "atomic_char_t"
	ATOMIC_INT8 = "atomic_int8_t" // @alias

	ATOMIC_SHORT = "atomic_short_t"
	ATOMIC_INT16 = "atomic_int16_t" // @alias

	ATOMIC_INT   = "atomic_int_t"
	ATOMIC_INT32 = "atomic_int32_t" // @alias

	ATOMIC_LONG  = "atomic_long_t"
	ATOMIC_INT64 = "atomic_long_t" // @alias

	ATOMIC_LLONG = "atomic_llong_t"

	ATOMIC_FLOAT = "atomic_float_t"

	ATOMIC_DOUBLE = "atomic_double_t"

	ATOMIC_PTR = "atomic_ptr_t"
)

type Interface struct {
	Funcs map[string]*ir.Func
	Types map[string]types.Type
}

var Instance *Interface
var once sync.Once

func NewInterface(mod *ir.Module) *Interface {
	t := &Interface{}
	once.Do(func() {
		t.Funcs = make(map[string]*ir.Func)
		t.Types = make(map[string]types.Type)
		t.RegisterTypes()
		t.registerFuncs(mod)
		Instance = t
	})
	return Instance
}
