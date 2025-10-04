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

func (t *IO) printf(typeHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) (typedef.Var, *ir.Block) {
	if len(args) == 0 {
		panic("printf requires at least one argument (format string)")
	}

	printfFn := t.ensurePrintf(module)

	// First arg must be a string
	formatVar := args[0]
	formatVal := formatVar.Load(block)
	if formatVal.Type() != types.I8Ptr {
		panic(fmt.Sprintf("printf: first argument must be string (i8*), got %s", formatVal.Type()))
	}

	callArgs := []value.Value{formatVal}

	// Remaining args
	for _, arg := range args[1:] {
		// Always load (so we pass value, not alloca slot)
		loaded := arg.Load(block)

		// printf is variadic, so no exact type match needed, but we should normalize:
		switch loaded.Type().(type) {
		case *types.IntType:
			// leave as-is (i32, i64, etc.)
			callArgs = append(callArgs, loaded)

		case *types.FloatType:
			// float must be promoted to double in varargs
			if loaded.Type() == types.Float {
				promoted := block.NewFPExt(loaded, types.Double)
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

	// Emit call
	result := block.NewCall(printfFn, callArgs...)

	// Wrap result in a Var (since printf returns int)
	return typeHandler.BuildVar(block, typedef.INT32, result), block
}
