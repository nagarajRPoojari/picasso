package types

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/ints"
)

type TypeHandler struct {
}

func NewTypeHandler() *TypeHandler {
	return &TypeHandler{}
}

func (t *TypeHandler) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["size"] = t.size
	funcs["type"] = t._type
	return funcs
}

func (t *TypeHandler) size(typeHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) typedef.Var {
	// assume args[0] holds the type or var we want sizeof
	typ := args[0].Type() // this should be `types.Type`

	// Create null pointer of given type
	nullPtr := constant.NewNull(types.NewPointer(typ))

	// GEP: move by 1 element
	gep := block.NewGetElementPtr(typ, nullPtr, constant.NewInt(types.I32, 1))

	// Convert pointer to int (i64)
	sizeVal := block.NewPtrToInt(gep, types.I64)
	slot := block.NewAlloca(types.I64)
	block.NewStore(sizeVal, slot)
	return &ints.Int32{
		NativeType: types.I64,
		Value:      slot,
	}
}

func (t *TypeHandler) _type(typeHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) typedef.Var {
	typ := args[0].NativeTypeString()

	strConst := constant.NewCharArrayFromString(typ + "\x00")
	global := module.NewGlobalDef("", strConst)

	gep := block.NewGetElementPtr(
		global.ContentType,
		global,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	return tf.NewString(block, gep)
}
