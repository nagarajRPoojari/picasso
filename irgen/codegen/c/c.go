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

	// Includes both standard C malloc and GC-tracked allocation.
	FUNC_MEMCPY  = "memcpy"
	ALIAS_MEMCPY = "memcpy"

	FUNC_RUNTIME_INIT  = "runtime_init"
	ALIAS_RUNTIME_INIT = "runtime_init"

	FUNC_RUNTIME_ERROR = "__public__runtime_error"

	FUNC_ALLOC  = "__public__alloc" // Garbage Collector tracked allocation
	ALIAS_ALLOC = "alloc"

	FUNC_ARRAY_ALLOC  = "__public__alloc_array"
	ALIAS_ARRAY_ALLOC = "alloc_array"

	TYPE_ARRAY   = "array"
	TYPE_RWMUTEX = "rwmutex"
	TYPE_MUTEX   = "mutex"

	// Runtime Type Names
	// These identifiers are used when declaring opaque or alias types in LLVM IR.
	TYPE_ATOMIC_BOOL = "atomic_bool"

	TYPE_ATOMIC_INT8  = "atomic_int8_t"
	TYPE_ATOMIC_INT16 = "atomic_int16_t"
	TYPE_ATOMIC_INT32 = "atomic_int32_t"
	TYPE_ATOMIC_INT64 = "atomic_int64_t"

	TYPE_ATOMIC_UINT8  = "atomic_uint8_t"
	TYPE_ATOMIC_UINT16 = "atomic_uint16_t"
	TYPE_ATOMIC_UINT32 = "atomic_uint32_t"
	TYPE_ATOMIC_UINT64 = "atomic_uint64_t"

	TYPE_ATOMIC_FLOAT16 = "atomic_float16_t" // 16-bit half
	TYPE_ATOMIC_FLOAT32 = "atomic_float32_t" // 32-bit single
	TYPE_ATOMIC_FLOAT64 = "atomic_float64_t" // 64-bit double

	TYPE_ATOMIC_PTR = "atomic_uintptr_t"
)

// Interface maintains a registry of available external functions and
// runtime types. It provides a centralized lookup table for the code
// generator to reference LLVM symbols.
type Interface struct {
	// Funcs maps symbol names to their corresponding LLVM IR function declarations.
	Funcs map[string]*ir.Func
	// Types maps type identifiers to their concrete LLVM IR type definitions.
	Types map[string]types.Type
	// constants
	Constants map[string]*ir.Global
}

// Instance is a global singleton providing access to the C runtime interface.
var Instance *Interface

// NewInterface initializes a new runtime registry for the given LLVM module.
// It populates the internal maps by registering all required external
// functions and built-in types.
func InitInterface(mod *ir.Module) *Interface {
	t := &Interface{}
	t.Funcs = make(map[string]*ir.Func)
	t.Types = make(map[string]types.Type)
	t.Constants = make(map[string]*ir.Global)

	t.registerTypes(mod)
	t.registerFuncs(mod)

	Instance = t
	return Instance
}
