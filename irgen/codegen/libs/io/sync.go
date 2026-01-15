package io

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/utils"

	function "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/func"
	typedef "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/floats"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/ints"
)

const (
	ALIAS_SPRINTF = "printf"
	ALIAS_SSCANF  = "scanf"
)

type SyncIO struct {
}

func NewSyncIO() *SyncIO {
	return &SyncIO{}
}

func (t *SyncIO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[ALIAS_SPRINTF] = t.sprintf
	funcs[ALIAS_SSCANF] = t.sscanf
	return funcs
}

func (t *SyncIO) sprintf(f *ir.Func, typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	castedArgs := []value.Value{args[0].Load(bh)}
	for _, arg := range args[1:] {
		switch arg.(type) {
		case *ints.Int8, *ints.Int16, *ints.Int32:
			res := typeHandler.ImplicitIntCast(bh, arg.Load(bh), types.I32)
			castedArgs = append(castedArgs, res)
		case *ints.Int64:
			castedArgs = append(castedArgs, arg.Load(bh))
		case *ints.UInt8, *ints.UInt16:
			res := typeHandler.ImplicitUnsignedIntCast(bh, arg.Load(bh), types.I32)
			castedArgs = append(castedArgs, res)
		case *ints.UInt32, *ints.UInt64:
			res := typeHandler.ImplicitUnsignedIntCast(bh, arg.Load(bh), types.I64)
			castedArgs = append(castedArgs, res)
		case *floats.Float16, *floats.Float32, *floats.Float64:
			res := typeHandler.ImplicitFloatCast(bh, arg.Load(bh), types.Double)
			castedArgs = append(castedArgs, res)
		default:
			castedArgs = append(castedArgs, arg.Load(bh))
		}
	}
	result := bh.N.NewCall(f, castedArgs...)
	return typeHandler.BuildVar(bh, typedef.NewType(utils.GetTypeString(result.Type())), result)
}

func (t *SyncIO) sscanf(f *ir.Func, typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	format := args[0].Load(bh)
	callArgs := []value.Value{format}

	for _, arg := range args[1:] {
		var slot value.Value
		switch arg.NativeTypeString() {
		case "string":
			slot = arg.Load(bh)
		default:
			slot = arg.Slot()
		}

		callArgs = append(callArgs, slot)
	}

	result := bh.N.NewCall(f, callArgs...)
	return typeHandler.BuildVar(bh, typedef.NewType(utils.GetTypeString(result.Type())), result)
}
