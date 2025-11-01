package utils

import (
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func BtoI(b bool) int64 {
	if b {
		return 1
	}
	return 0
}

func AtomicStore(block *ir.Block, dst, val value.Value, align int) {
	inst := block.NewStore(val, dst)
	inst.Atomic = true
	inst.Ordering = enum.AtomicOrderingSequentiallyConsistent
	inst.Align = ir.Align(align)
}

func AtomicLoad(block *ir.Block, typ types.Type, src value.Value, align int) value.Value {
	inst := block.NewLoad(typ, src)
	inst.Atomic = true
	inst.Ordering = enum.AtomicOrderingSequentiallyConsistent
	inst.Align = ir.Align(align)
	return inst
}
