package sync

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/compiler/c"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/utils"

	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/libutils"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	bc "github.com/nagarajRPoojari/x-lang/compiler/type/block"
)

type Sync struct {
}

func NewSync() *Sync {
	return &Sync{}

}

func (t *Sync) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.PRINTF] = t.printf
	funcs[c.ATOMIC_STORE] = t.atomicStore
	funcs[c.ATOMIC_LOAD] = t.atomicLoad
	funcs[c.ATOMIC_ADD] = t.atomicAdd
	funcs[c.ATOMIC_SUB] = t.atomicSub
	return funcs
}

func (t *Sync) printf(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	printfFn := c.NewInterface(module).Funcs[c.PRINTF]
	return libutils.CallCFunc(typeHandler, printfFn, bh, args)
}

func (t *Sync) atomicStore(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.NewInterface(module).Funcs[fmt.Sprintf("atomic_store_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicLoad(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.NewInterface(module).Funcs[fmt.Sprintf("atomic_load_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicAdd(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.NewInterface(module).Funcs[fmt.Sprintf("atomic_add_%s", extractBetween(dest, "_", "_"))]
	return libutils.CallCFunc(typeHandler, method, bh, args)
}

func (t *Sync) atomicSub(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	dest := utils.GetTypeString(args[0].Type())
	method := c.NewInterface(module).Funcs[fmt.Sprintf("atomic_sub_%s", extractBetween(dest, "_", "_"))]
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
