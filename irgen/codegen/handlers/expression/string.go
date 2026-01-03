package expression

import (
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// ProcessStringLiteral generates LLVM IR for constant string expressions.
// Unlike simple primitives, strings in Niyama are treated as heap-allocated
// objects to support mutability and lifetime management beyond the local scope.
//
// Technical Logic:
//   - Global Definition: Stores the raw string literal as an internal constant
//     in the LLVM module's data section, ensuring a null terminator (\x00) is added.
//   - Heap Allocation: Invokes the runtime's allocation function (malloc) to
//     reserve space in the heap equivalent to the string length plus the terminator.
//   - Memory Migration: Uses a 'memcpy' operation to copy the read-only data from
//     the global segment into the newly allocated heap memory.
//   - Pointer Wrapping: Returns a Niyama string variable container that points
//     to the addressable heap memory.
func (t *ExpressionHandler) ProcessStringLiteral(bh *bc.BlockHolder, ex ast.StringExpression) tf.Var {
	formatStr := ex.Value

	strConst := constant.NewCharArrayFromString(formatStr + "\x00")
	global := t.st.Module.NewGlobalDef("", strConst)

	// string vars are heap allocated & stored as pointer to global constants.
	// set linkage to internal to avoid name collision with linked IR modules.
	global.Linkage = enum.LinkageInternal
	gep := bh.N.NewGetElementPtr(
		global.ContentType,
		global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	return tf.NewString(bh, gep)
}
