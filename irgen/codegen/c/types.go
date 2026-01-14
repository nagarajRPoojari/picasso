package c

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

// registerTypes initializes the type system for the LLVM module. It ensures
// that complex structures like arrays and synchronization primitives are
// defined before they are referenced by instructions.
func (t *Interface) registerTypes(mod *ir.Module) {
	t.initArrayTypes(mod)
	t.initAtomicTypes(mod)
}

// initArrayTypes defines the internal representation of the Niyama array.
// It creates a struct containing a generic data pointer, a shape descriptor,
// the flat length, and the dimensional rank. This definition is registered
// as a named type in the LLVM module for cross-function consistency
func (t *Interface) initArrayTypes(_ *ir.Module) {
	t.Types[TYPE_ARRAY] = types.NewStruct(
		types.NewPointer(types.I8),  // data
		types.NewPointer(types.I64), // shape (i64*)
		types.I64,                   // length
		types.I64,                   // rank
	)
}

// initAtomicTypes wraps fundamental scalar types in LLVM structures to
// represent atomic variables. Wrapping these in structs provides a clear
// distinction between standard and thread-safe variables during type
// checking and IR lowering.
func (t *Interface) initAtomicTypes(_ *ir.Module) {
	// Boolean
	t.Types[TYPE_ATOMIC_BOOL] = types.NewStruct(types.I1)
	t.Types["i1"] = types.NewStruct(types.I1)

	// Exact-width signed integers
	t.Types[TYPE_ATOMIC_INT8] = types.NewStruct(types.I8)
	t.Types["i8"] = types.NewStruct(types.I8)
	t.Types[TYPE_ATOMIC_INT16] = types.NewStruct(types.I16)
	t.Types["i16"] = types.NewStruct(types.I16)
	t.Types[TYPE_ATOMIC_INT32] = types.NewStruct(types.I32)
	t.Types["i32"] = types.NewStruct(types.I32)
	t.Types[TYPE_ATOMIC_INT64] = types.NewStruct(types.I64)
	t.Types["i64"] = types.NewStruct(types.I64)

	// Exact-width unsigned integers
	t.Types[TYPE_ATOMIC_UINT8] = types.NewStruct(types.I8)
	t.Types[TYPE_ATOMIC_UINT16] = types.NewStruct(types.I16)
	t.Types[TYPE_ATOMIC_UINT32] = types.NewStruct(types.I32)
	t.Types[TYPE_ATOMIC_UINT64] = types.NewStruct(types.I64)

	// Floating point
	t.Types[TYPE_ATOMIC_FLOAT16] = types.NewStruct(types.Half)   // 16-bit
	t.Types[TYPE_ATOMIC_FLOAT32] = types.NewStruct(types.Float)  // 32-bit
	t.Types[TYPE_ATOMIC_FLOAT64] = types.NewStruct(types.Double) // 64-bit

	// Pointer (atomic_uintptr_t)
	t.Types[TYPE_ATOMIC_PTR] = types.NewStruct(types.NewPointer(types.I8))

	t.Types[TYPE_RWMUTEX] = types.NewStruct()

	t.Types[TYPE_MUTEX] = types.NewStruct()
}
