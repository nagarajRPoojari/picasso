package rterr

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

const (
	RUNTIME_ERR = "__public__runtime_error"
)

type ErrorHandler struct {
	err *ir.Func
}

var Instance *ErrorHandler

func NewErrorHandler(mod *ir.Module) *ErrorHandler {
	Instance = &ErrorHandler{
		err: mod.NewFunc(RUNTIME_ERR, types.Void, ir.NewParam("msg", types.I8Ptr)),
	}
	return Instance
}

func (t *ErrorHandler) RaiseRTError(block *ir.Block, msg string) {
	m := block.Parent.Parent
	msgGlobal := m.NewGlobalDef("", constant.NewCharArrayFromString(msg))
	msgGlobal.Immutable = true
	msgGlobal.Linkage = enum.LinkagePrivate
	msgPtr := block.NewGetElementPtr(msgGlobal.Type(), msgGlobal,
		constant.NewInt(types.I32, 0),
	)
	block.NewCall(t.err, msgPtr)
	block.NewUnreachable()
}
