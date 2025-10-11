package c

import "github.com/llir/llvm/ir/types"

var Array = types.NewStruct(
	types.I64,                   // length
	types.NewPointer(types.I8),  // data
	types.NewPointer(types.I64), // shape (i64*)
	types.I64,                   // rank
)
