package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	function "github.com/nagarajRPoojari/x-lang/generator/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	typedef "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
	"github.com/nagarajRPoojari/x-lang/generator/type/primitives/ints"
)

type TypeHandler struct {
}

func NewTypeHandler() *TypeHandler {
	return &TypeHandler{}
}

func (t *TypeHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["size"] = t.Size
	funcs["type"] = t.TypeOf
	return funcs
}

func (t *TypeHandler) Size(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	// assume args[0] holds the type or var we want sizeof
	typ := args[0].Type() // this should be `types.Type`

	// Create null pointer of given type
	nullPtr := constant.NewNull(types.NewPointer(typ))

	// GEP: move by 1 element
	gep := bh.N.NewGetElementPtr(typ, nullPtr, constant.NewInt(types.I32, 1))

	// Convert pointer to int (i64)
	sizeVal := bh.N.NewPtrToInt(gep, types.I64)
	slot := bh.V.NewAlloca(types.I64)
	bh.N.NewStore(sizeVal, slot)
	return &ints.Int32{
		NativeType: types.I64,
		Value:      slot,
	}
}

func (t *TypeHandler) TypeOf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	typ := args[0].NativeTypeString()

	strConst := constant.NewCharArrayFromString(typ + "\x00")
	global := module.NewGlobalDef("", strConst)

	gep := bh.N.NewGetElementPtr(
		global.ContentType,
		global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	return tf.NewString(bh, gep)
}
