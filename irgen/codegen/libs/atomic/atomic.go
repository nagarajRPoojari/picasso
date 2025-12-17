package sync

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"

	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs/libutils"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

type Sync struct {
}

func NewSync() *Sync {
	return &Sync{}

}

func (t *Sync) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.ALIAS_ATOMIC_STORE] = t.atomicStore
	funcs[c.ALIAS_ATOMIC_LOAD] = t.atomicLoad
	funcs[c.ALIAS_ATOMIC_ADD] = t.atomicAdd
	funcs[c.ALIAS_ATOMIC_SUB] = t.atomicSub
	return funcs
}

func (t *Sync) atomicStore(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.Instance.Funcs[fmt.Sprintf("__public__atomic_store_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicLoad(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.Instance.Funcs[fmt.Sprintf("__public__atomic_load_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicAdd(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.Instance.Funcs[fmt.Sprintf("__public__atomic_add_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicSub(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.Instance.Funcs[fmt.Sprintf("__public__atomic_sub_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func extractBetween(s, start, end string) string {
	startIdx := strings.Index(s, start)
	if startIdx == -1 {
		return ""
	}
	startIdx += len(start)

	endIdx := strings.Index(s[startIdx:], end)
	if endIdx == -1 {
		return ""
	}
	return s[startIdx : startIdx+endIdx]
}
