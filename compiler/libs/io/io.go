package io

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
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

func (t *IO) printf(typeHandler *tf.TypeHandler, module *ir.Module, block *ir.Block, args []typedef.Var) typedef.Var {
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
				callArgs = append(callArgs, loaded) // already double/half/etc.
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
	return typeHandler.BuildVar(block, typedef.INT32, result)
}

func (t *IO) printfUtil(module *ir.Module, block *ir.Block, format string, args ...interface{}) {
	printf := t.ensurePrintf(module)

	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	if !strings.HasSuffix(format, "\x00") {
		format += "\x00"
	}

	strConst := constant.NewCharArrayFromString(format)
	x := uuid.NewString()
	global := module.NewGlobalDef(fmt.Sprintf(".str%s", x), strConst)
	global.Linkage = enum.LinkagePrivate

	zero := constant.NewInt(types.I64, 0)
	ptr := block.NewGetElementPtr(strConst.Typ, global, zero, zero)

	var callArgs []value.Value
	callArgs = append(callArgs, ptr)

	for _, a := range args {
		switch v := a.(type) {
		case value.Value:
			// Already an LLVM value (float, int, etc.)
			callArgs = append(callArgs, v)

		case string:
			// Turn Go string into global char array + pointer
			str := v + "\x00"
			strConst := constant.NewCharArrayFromString(str)

			x := uuid.NewString()
			g := module.NewGlobalDef(fmt.Sprintf(".str%s", x), strConst)
			g.Linkage = enum.LinkagePrivate

			zero := constant.NewInt(types.I64, 0)
			p := block.NewGetElementPtr(strConst.Typ, g, zero, zero)
			callArgs = append(callArgs, p)

		case int:
			callArgs = append(callArgs, constant.NewInt(types.I32, int64(v)))

		case float64:
			callArgs = append(callArgs, constant.NewFloat(types.Double, v))

		case float32:
			callArgs = append(callArgs, constant.NewFloat(types.Float, float64(v)))

		default:
			panic(fmt.Sprintf("unsupported print arg type: %T", v))
		}
	}

	block.NewCall(printf, callArgs...)
}
