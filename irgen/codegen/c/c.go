/*
Package c provides the Foreign Function Interface (FFI) and Runtime Bridge for
compiler. Exposes set of functions & types from c code.
Notes:
  - it is expected to link corresponding c binaries during runtime
  - functions are accessed through libs defined in libs/ directory
  - __public__ identifier indicates c functions name while corresponding alias
    used as exposed name
*/
package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// Foreign Function and Runtime identifiers.
// These constants define the mapping between internal function
// aliases and the actual symbol names in the C runtime or standard library.
const (
	// Debug and Diagnostic Utilities
	__UTILS__FUNC_DEBUG_ARRAY_INFO = "__public__debug_array_info"

	// Concurrency and Threading
	FUNC_THREAD      = "thread"
	ALIAS_THREAD     = "thread"
	FUNC_SELF_YIELD  = "self_yield"
	ALIAS_SELF_YIELD = "self_yield"

	// Blocking I/O Calls
	// These map directly to the underlying OS file descriptors.
	ALIAS_PRINTF  = "printf"
	ALIAS_FPRINTF = "fprintf"
	ALIAS_SCANF   = "scanf"
	ALIAS_FSCANF  = "fscanf"
	ALIAS_FREAD   = "fread"
	ALIAS_FWRITE  = "fwrite"

	FUNC_FOPEN   = "fopen"
	ALIAS_FOPEN  = "fopen"
	FUNC_FCLOSE  = "fclose"
	ALIAS_FCLOSE = "fclose"
	FUNC_FFLUSH  = "fflush"
	ALIAS_FFLUSH = "fflush"
	FUNC_FSEEK   = "fseek"
	ALIAS_FSEEK  = "fseek"
	FUNC_FPUTS   = "fputs"
	ALIAS_FPUTS  = "fputs"
	FUNC_FGETS   = "fgets"
	ALIAS_FGETS  = "fgets"

	// Formatted I/O and Buffer operations
	FUNC_FSCANF  = "fscanf"
	FUNC_FPRINTF = "fprintf"
	FUNC_SPRINTF = "__public__sprintf"
	FUNC_SSCAN   = "__public__sscan"
	FUNC_SFREAD  = "__public__sfread"
	FUNC_SFWRITE = "__public__sfwrite"

	// Non-blocking I/O (Asynchronous Wrapper Calls)
	FUNC_APRINTF = "__public__aprintf"
	FUNC_ASCAN   = "__public__ascan"
	FUNC_AFREAD  = "__public__afread"
	FUNC_AFWRITE = "__public__afwrite"

	// Memory Management
	// Includes both standard C malloc and GC-tracked allocation.
	FUNC_MALLOC   = "malloc"
	ALIAS_MALLOC  = "malloc"
	FUNC_MEMCPY   = "memcpy"
	ALIAS_MEMCPY  = "memcpy"
	FUNC_MEMSET   = "memset"
	ALIAS_MEMSET  = "memset"
	FUNC_MEMMOVE  = "memmove"
	ALIAS_MEMMOVE = "memmove"

	FUNC_RUNTIME_INIT  = "runtime_init"
	ALIAS_RUNTIME_INIT = "runtime_init"

	FUNC_ALLOC  = "__public__alloc" // Garbage Collector tracked allocation
	ALIAS_ALLOC = "alloc"

	FUNC_ARRAY_ALLOC  = "__public__alloc_array"
	ALIAS_ARRAY_ALLOC = "alloc_array"

	TYPE_ARRAY = "array"

	// String and Container Operations
	FUNC_STRLEN   = "strlen"
	ALIAS_STRLEN  = "strlen"
	FUNC_FORMAT   = "format"
	ALIAS_FORMAT  = "format"
	FUNC_LEN      = "len"
	ALIAS_LEN     = "len"
	FUNC_COMPARE  = "compare"
	ALIAS_COMPARE = "compare"
	FUNC_STRCMP   = "strcmp"
	ALIAS_STRCMP  = "strcmp"

	// Process Control
	FUNC_EXIT  = "exit"
	ALIAS_EXIT = "exit"
	FUNC_HASH  = "hash"
	ALIAS_HASH = "hash"

	// Atomic Operations and Synchronization
	// These constants map to thread-safe primitives for various bit-widths.
	ALIAS_ATOMIC_STORE = "atomic_store"
	ALIAS_ATOMIC_LOAD  = "atomic_load"
	ALIAS_ATOMIC_ADD   = "atomic_add"
	ALIAS_ATOMIC_SUB   = "atomic_sub"

	ATOMIC_STORE_BOOL = "__public__atomic_store_bool"
	ATOMIC_LOAD_BOOL  = "__public__atomic_load_bool"

	// 8-bit Atomics
	ATOMIC_STORE_CHAR = "__public__atomic_store_char"
	ATOMIC_STORE_INT8 = "__public__atomic_store_int8"
	ATOMIC_LOAD_CHAR  = "__public__atomic_load_char"
	ATOMIC_LOAD_INT8  = "__public__atomic_load_int8"
	ATOMIC_ADD_CHAR   = "__public__atomic_add_char"
	ATOMIC_ADD_INT8   = "__public__atomic_add_int8"
	ATOMIC_SUB_CHAR   = "__public__atomic_sub_char"
	ATOMIC_SUB_INT8   = "__public__atomic_sub_int8"

	// 16-bit Atomics
	ATOMIC_STORE_SHORT = "__public__atomic_store_short"
	ATOMIC_STORE_INT16 = "__public__atomic_store_int16"
	ATOMIC_LOAD_SHORT  = "__public__atomic_load_short"
	ATOMIC_LOAD_INT16  = "__public__atomic_load_int16"
	ATOMIC_ADD_SHORT   = "__public__atomic_add_short"
	ATOMIC_ADD_INT16   = "__public__atomic_add_int16"
	ATOMIC_SUB_SHORT   = "__public__atomic_sub_short"
	ATOMIC_SUB_INT16   = "__public__atomic_sub_int16"

	// 32-bit Atomics
	ATOMIC_STORE_INT   = "__public__atomic_store_int"
	ATOMIC_STORE_INT32 = "__public__atomic_store_int32"
	ATOMIC_LOAD_INT    = "__public__atomic_load_int"
	ATOMIC_LOAD_INT32  = "__public__atomic_load_int32"
	ATOMIC_ADD_INT     = "__public__atomic_add_int"
	ATOMIC_ADD_INT32   = "__public__atomic_add_int32"
	ATOMIC_SUB_INT     = "__public__atomic_sub_int"
	ATOMIC_SUB_INT32   = "__public__atomic_sub_int32"

	// 64-bit Atomics
	ATOMIC_STORE_LONG  = "__public__atomic_store_long"
	ATOMIC_STORE_INT64 = "__public__atomic_store_int64"
	ATOMIC_LOAD_LONG   = "__public__atomic_load_long"
	ATOMIC_LOAD_INT64  = "__public__atomic_load_int64"
	ATOMIC_ADD_LONG    = "__public__atomic_add_long"
	ATOMIC_ADD_INT64   = "__public__atomic_add_int64"
	ATOMIC_SUB_LONG    = "__public__atomic_sub_long"
	ATOMIC_SUB_INT64   = "__public__atomic_sub_int64"

	// Pointer and Floating Point Atomics
	ATOMIC_STORE_LLONG  = "__public__atomic_store_llong"
	ATOMIC_LOAD_LLONG   = "__public__atomic_load_llong"
	ATOMIC_ADD_LLONG    = "__public__atomic_add_llong"
	ATOMIC_SUB_LLONG    = "__public__atomic_sub_llong"
	ATOMIC_STORE_FLOAT  = "__public__atomic_store_float"
	ATOMIC_LOAD_FLOAT   = "__public__atomic_load_float"
	ATOMIC_STORE_DOUBLE = "__public__atomic_store_double"
	ATOMIC_LOAD_DOUBLE  = "__public__atomic_load_double"
	ATOMIC_STORE_PTR    = "__public__atomic_store_ptr"
	ATOMIC_LOAD_PTR     = "__public__atomic_load_ptr"

	// Runtime Type Names
	// These identifiers are used when declaring opaque or alias types in LLVM IR.
	TYPE_ATOMIC_BOOL   = "atomic_bool_t"
	TYPE_ATOMIC_CHAR   = "atomic_char_t"
	TYPE_ATOMIC_INT8   = "atomic_int8_t"
	TYPE_ATOMIC_SHORT  = "atomic_short_t"
	TYPE_ATOMIC_INT16  = "atomic_int16_t"
	TYPE_ATOMIC_INT    = "atomic_int_t"
	TYPE_ATOMIC_INT32  = "atomic_int32_t"
	TYPE_ATOMIC_LONG   = "atomic_long_t"
	TYPE_ATOMIC_INT64  = "atomic_long_t"
	TYPE_ATOMIC_LLONG  = "atomic_llong_t"
	TYPE_ATOMIC_FLOAT  = "atomic_float_t"
	TYPE_ATOMIC_DOUBLE = "atomic_double_t"
	TYPE_ATOMIC_PTR    = "atomic_ptr_t"
)

// Interface maintains a registry of available external functions and
// runtime types. It provides a centralized lookup table for the code
// generator to reference LLVM symbols.
type Interface struct {
	// Funcs maps symbol names to their corresponding LLVM IR function declarations.
	Funcs map[string]*ir.Func
	// Types maps type identifiers to their concrete LLVM IR type definitions.
	Types map[string]types.Type
}

// Instance is a global singleton providing access to the C runtime interface.
var Instance *Interface

// NewInterface initializes a new runtime registry for the given LLVM module.
// It populates the internal maps by registering all required external
// functions and built-in types.
func NewInterface(mod *ir.Module) *Interface {
	t := &Interface{}
	t.Funcs = make(map[string]*ir.Func)
	t.Types = make(map[string]types.Type)

	t.registerTypes(mod)
	t.registerFuncs(mod)

	Instance = t
	return Instance
}
