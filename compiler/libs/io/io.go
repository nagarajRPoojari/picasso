package io

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/compiler/c"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"

	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/libutils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	typedef "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

type IO struct {
}

func NewIO() *IO {
	return &IO{}

}

func (t *IO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.PRINTF] = t.printf
	funcs[c.SCANF] = t.scanf
	funcs[c.FOPEN] = t.fopen
	funcs[c.FCLOSE] = t.fclose
	funcs[c.FPRINTF] = t.fprintf
	funcs[c.FSCANF] = t.fscanf
	funcs[c.FPUTS] = t.fputs
	funcs[c.FGETS] = t.fgets
	funcs[c.FFLUSH] = t.fflush
	funcs[c.FSEEK] = t.fseek
	return funcs
}

func (t *IO) printf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	printfFn := c.NewInterface(module).Funcs[c.PRINTF]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *IO) scanf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanfFunc := c.NewInterface(module).Funcs[c.SCANF]

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

	result := bh.N.NewCall(scanfFunc, callArgs...)
	return typeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(result.Type())), result)
}

func (t *IO) fopen(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fopenFunc := c.NewInterface(module).Funcs[c.FOPEN]
	return libutils.CallCFunc(typeHandler, fopenFunc, bh, args)
}

func (t *IO) fprintf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fprintfFn := c.NewInterface(module).Funcs[c.FPRINTF]
	return libutils.CallCFunc(typeHandler, fprintfFn, bh, args)
}

func (t *IO) fscanf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fscanfFunc := c.NewInterface(module).Funcs[c.FSCANF]

	file := args[0].Load(bh)
	format := args[1].Load(bh)
	callArgs := []value.Value{file, format}

	for _, arg := range args[2:] {
		var slot value.Value
		switch arg.NativeTypeString() {
		case "string":
			slot = arg.Load(bh)
		default:
			slot = arg.Slot()
		}

		callArgs = append(callArgs, slot)
	}

	result := bh.N.NewCall(fscanfFunc, callArgs...)
	return typeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(result.Type())), result)
}

func (t *IO) fclose(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fscloseFunc := c.NewInterface(module).Funcs[c.FCLOSE]
	return libutils.CallCFunc(typeHandler, fscloseFunc, bh, args)
}

func (t *IO) fputs(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fputsFunc := c.NewInterface(module).Funcs[c.FPUTS]
	return libutils.CallCFunc(typeHandler, fputsFunc, bh, args)
}

func (t *IO) fgets(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fgetsFunc := c.NewInterface(module).Funcs[c.FGETS]
	return libutils.CallCFunc(typeHandler, fgetsFunc, bh, args)
}

func (t *IO) fflush(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fflushFunc := c.NewInterface(module).Funcs[c.FFLUSH]
	return libutils.CallCFunc(typeHandler, fflushFunc, bh, args)
}

func (t *IO) fseek(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fseekFunc := c.NewInterface(module).Funcs[c.FSEEK]
	return libutils.CallCFunc(typeHandler, fseekFunc, bh, args)
}
