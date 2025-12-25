package net

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs/libutils"
	rterr "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/private/runtime"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

type NetIO struct {
}

func NewNetIO() *NetIO {
	return &NetIO{}
}

func (t *NetIO) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.ALIAS_NET_LISTEN] = t.listen
	funcs[c.ALIAS_NET_ACCEPT] = t.accept
	funcs[c.ALIAS_NET_READ] = t.read
	funcs[c.ALIAS_NET_WRITE] = t.write
	return funcs
}

func (t *NetIO) listen(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	listen := c.Instance.Funcs[c.FUNC_NET_LISTEN]
	return libutils.CallCFunc(typeHandler, listen, bh, args)
}

func (t *NetIO) accept(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	accept := c.Instance.Funcs[c.FUNC_NET_ACCEPT]
	return libutils.CallCFunc(typeHandler, accept, bh, args)
}

func (t *NetIO) read(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanFn := c.Instance.Funcs[c.FUNC_NET_READ]

	// args: fd, buf, length
	dest := args[1]
	if utils.GetTypeString(dest.Type()) != "array" {
		errorutils.Abort(errorutils.ParamsError, "i8*", dest.Type())
	}
	size := typeHandler.ImplicitTypeCast(bh, utils.GetTypeString(scanFn.Sig.Params[2]), args[2].Load(bh))

	CheckIntCond(bh, dest.(*typedef.Array).Len(bh), size, enum.IPredSGE, "buffer overflow")

	return libutils.CallCFunc(typeHandler, scanFn, bh, args)
}

// @todo: fix
func (t *NetIO) write(typeHandler *typedef.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	scanFn := c.Instance.Funcs[c.FUNC_NET_WRITE]

	// dest := args[1]
	// if utils.GetTypeString(dest.Type()) != "array" {
	// 	errorutils.Abort(errorutils.ParamsError, "i8*", dest.Type())
	// }
	// size := typeHandler.ImplicitTypeCast(bh, utils.GetTypeString(scanFn.Sig.Params[2]), args[2].Load(bh))

	// CheckIntCond(bh, dest.(*typedef.Array).Len(bh), size, enum.IPredSLT, "buffer underflow")

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
