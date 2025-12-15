package c

import "github.com/llir/llvm/ir/types"

func (t *Interface) RegisterTypes() {
	t.initArrayTypes()
	t.initAtomicTypes()
}

func (t *Interface) initArrayTypes() {
	t.Types[TYPE_ARRAY] = types.NewStruct(
		types.I64,                   // length
		types.NewPointer(types.I8),  // data
		types.NewPointer(types.I64), // shape (i64*)
		types.I64,                   // rank
	)
}

func (t *Interface) initAtomicTypes() {
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
