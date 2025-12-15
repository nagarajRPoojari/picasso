package io

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/frontend/codegen/error"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/handlers/utils"
	function "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/libutils"
	rterr "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/private/runtime"
	typedef "github.com/nagarajRPoojari/niyama/frontend/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/frontend/codegen/type/block"
)

type AsyncIO struct {
}

func NewAsyncIO() *AsyncIO {
	return &AsyncIO{}
}

func (t *AsyncIO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.ALIAS_PRINTF] = t.aprintf
	funcs[c.ALIAS_SCANF] = t.ascan
	funcs[c.ALIAS_FREAD] = t.afread
	funcs[c.ALIAS_FWRITE] = t.awrite
	return funcs
}

func (t *AsyncIO) aprintf(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	printfFn := c.NewInterface(module).Funcs[c.FUNC_APRINTF]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *AsyncIO) ascan(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_ASCAN]
	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}

func (t *AsyncIO) afread(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_AFREAD]

	dest := args[1]
	if utils.GetTypeString(dest.Type()) != "array" {
		errorutils.Abort(errorutils.ParamsError, "i8*", dest.Type())
	}
	size := typeHandler.ImplicitTypeCast(bh, utils.GetTypeString(scanFn.Sig.Params[2]), args[2].Load(bh))

	CheckIntCond(bh, dest.(*typedef.Array).Len(bh), size, enum.IPredSGE, "buffer overflow")

	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}

// @todo: fix
func (t *AsyncIO) awrite(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanFn := c.NewInterface(module).Funcs[c.FUNC_AFWRITE]

	dest := args[1]
	if utils.GetTypeString(dest.Type()) != "array" {
		errorutils.Abort(errorutils.ParamsError, "i8*", dest.Type())
	}
	size := typeHandler.ImplicitTypeCast(bh, utils.GetTypeString(scanFn.Sig.Params[2]), args[2].Load(bh))

	CheckIntCond(bh, dest.(*typedef.Array).Len(bh), size, enum.IPredSLT, "buffer underflow")

	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}

func CheckIntCond(
	block *bc.BlockHolder,
	v1, v2 value.Value,
	pred enum.IPred,
	errMsg string,
) {
	b := block.N

	passBlk := b.Parent.NewBlock("")
	failBlk := b.Parent.NewBlock("")

	cond := b.NewICmp(pred, v1, v2)

	b.NewCondBr(cond, passBlk, failBlk)

	rterr.Instance.RaiseRTError(failBlk, errMsg)

	block.Update(block.V, passBlk)
}
