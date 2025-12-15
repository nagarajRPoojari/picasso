package io

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/c"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/utils"

	function "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/libutils"
	typedef "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

type SyncIO struct {
}

func NewSyncIO() *SyncIO {
	return &SyncIO{}
}

func (t *SyncIO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.ALIAS_PRINTF] = t.sprintf
	funcs[c.ALIAS_SCANF] = t.sscanf
	funcs[c.ALIAS_FREAD] = t.sfread
	funcs[c.ALIAS_FWRITE] = t.sfwrite

	funcs[c.ALIAS_FPRINTF] = t.sfprintf
	funcs[c.ALIAS_FSCANF] = t.sfscanf
	funcs[c.ALIAS_FPUTS] = t.sfputs
	funcs[c.ALIAS_FGETS] = t.sfgets

	funcs[c.ALIAS_FOPEN] = t.fopen
	funcs[c.ALIAS_FCLOSE] = t.fclose
	funcs[c.ALIAS_FFLUSH] = t.fflush
	funcs[c.ALIAS_FSEEK] = t.fseek
	return funcs
}

func (t *SyncIO) sprintf(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	printfFn := c.NewInterface(module).Funcs[c.FUNC_SPRINTF]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *SyncIO) sscanf(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanfFunc := c.NewInterface(module).Funcs[c.FUNC_SSCAN]

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
	return typeHandler.BuildVar(bh, typedef.NewType(utils.GetTypeString(result.Type())), result)
}

func (t *SyncIO) sfread(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	printfFn := c.NewInterface(module).Funcs[c.FUNC_SFREAD]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *SyncIO) sfwrite(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	printfFn := c.NewInterface(module).Funcs[c.FUNC_SFWRITE]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *SyncIO) fopen(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fopenFunc := c.NewInterface(module).Funcs[c.FUNC_FOPEN]
	return libutils.CallCFunc(typeHandler, fopenFunc, bh, args)
}

func (t *SyncIO) sfprintf(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fprintfFn := c.NewInterface(module).Funcs[c.FUNC_FPRINTF]
	return libutils.CallCFunc(typeHandler, fprintfFn, bh, args)
}

func (t *SyncIO) sfscanf(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fscanfFunc := c.NewInterface(module).Funcs[c.FUNC_FSCANF]

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
	return typeHandler.BuildVar(bh, typedef.NewType(utils.GetTypeString(result.Type())), result)
}

func (t *SyncIO) fclose(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fscloseFunc := c.NewInterface(module).Funcs[c.FUNC_FCLOSE]
	return libutils.CallCFunc(typeHandler, fscloseFunc, bh, args)
}

func (t *SyncIO) sfputs(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fputsFunc := c.NewInterface(module).Funcs[c.FUNC_FPUTS]
	return libutils.CallCFunc(typeHandler, fputsFunc, bh, args)
}

func (t *SyncIO) sfgets(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fgetsFunc := c.NewInterface(module).Funcs[c.FUNC_FGETS]
	return libutils.CallCFunc(typeHandler, fgetsFunc, bh, args)
}

func (t *SyncIO) fflush(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fflushFunc := c.NewInterface(module).Funcs[c.FUNC_FFLUSH]
	return libutils.CallCFunc(typeHandler, fflushFunc, bh, args)
}

func (t *SyncIO) fseek(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	fseekFunc := c.NewInterface(module).Funcs[c.FUNC_FSEEK]
	return libutils.CallCFunc(typeHandler, fseekFunc, bh, args)
}
