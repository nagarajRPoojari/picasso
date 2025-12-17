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
func (t *Interface) initArrayTypes(mod *ir.Module) {
	t.Types[TYPE_ARRAY] = types.NewStruct(
		types.NewPointer(types.I8),  // data
		types.NewPointer(types.I64), // shape (i64*)
		types.I64,                   // length
		types.I64,                   // rank
	)

	mod.NewTypeDef(TYPE_ARRAY, t.Types[TYPE_ARRAY])
}

// initAtomicTypes wraps fundamental scalar types in LLVM structures to
// represent atomic variables. Wrapping these in structs provides a clear
// distinction between standard and thread-safe variables during type
// checking and IR lowering.
func (t *Interface) initAtomicTypes(mod *ir.Module) {
	t.Types[TYPE_ATOMIC_BOOL] = types.NewStruct(types.I1)

	t.Types[TYPE_ATOMIC_CHAR] = types.NewStruct(types.I8)
	t.Types[TYPE_ATOMIC_INT8] = types.NewStruct(types.I8)

	t.Types[TYPE_ATOMIC_SHORT] = types.NewStruct(types.I16)
	t.Types[TYPE_ATOMIC_INT16] = types.NewStruct(types.I16)

	t.Types[TYPE_ATOMIC_INT] = types.NewStruct(types.I32)
	t.Types[TYPE_ATOMIC_INT32] = types.NewStruct(types.I32)

	t.Types[TYPE_ATOMIC_LONG] = types.NewStruct(types.I64)
	t.Types[TYPE_ATOMIC_INT64] = types.NewStruct(types.I64)

	t.Types[TYPE_ATOMIC_LLONG] = types.NewStruct(types.I64)

	t.Types[TYPE_ATOMIC_FLOAT] = types.NewStruct(types.Float)

	t.Types[TYPE_ATOMIC_DOUBLE] = types.NewStruct(types.Double)

	t.Types[TYPE_ATOMIC_PTR] = types.NewStruct(types.NewPointer(types.I8))
}
