package io

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
)

type IO struct {
}

func NewIO() *IO {
	return &IO{}

}

func (t *IO) ensurePrintf(module *ir.Module) *ir.Func {
	for _, f := range module.Funcs {
		if f.Name() == "printf" {
			return f
		}
	}
	printf := module.NewFunc("printf", types.I32, ir.NewParam("", types.I8Ptr))
	printf.Sig.Variadic = true
	return printf
}

func (t *IO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs["printf"] = t.printf
	return funcs
}

func (t *IO) printf(typeHandler *tf.TypeHandler, module *ir.Module, bh tf.BlockHolder, args []typedef.Var) (typedef.Var, tf.BlockHolder) {
	if len(args) == 0 {
		panic("printf requires at least one argument (format string)")
	}

	printfFn := t.ensurePrintf(module)

	formatVar := args[0]
	formatVal := formatVar.Load(bh.N)
	if formatVal.Type() != types.I8Ptr {
		panic(fmt.Sprintf("printf: first argument must be string (i8*), got %s", formatVal.Type()))
	}

	callArgs := []value.Value{formatVal}

	for _, arg := range args[1:] {
		loaded := arg.Load(bh.N)

		switch loaded.Type().(type) {
		case *types.IntType:
			callArgs = append(callArgs, loaded)

		case *types.FloatType:
			if loaded.Type() == types.Float {
				promoted := bh.N.NewFPExt(loaded, types.Double)
				callArgs = append(callArgs, promoted)
			} else {
				callArgs = append(callArgs, loaded)
			}

		case *types.PointerType:
			callArgs = append(callArgs, loaded)

		default:
			panic(fmt.Sprintf("printf: unsupported argument type %s", loaded.Type()))
		}
	}

	result := bh.N.NewCall(printfFn, callArgs...)
	return typeHandler.BuildVar(bh, tf.NewType(typedef.INT32), result), bh
}
