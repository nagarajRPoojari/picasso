package libutils

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/utils"
	typedef "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

// @todo: validate params & raise error
func CallCFunc(typeHandler *typedef.TypeHandler, f *ir.Func, bh *bc.BlockHolder, args []typedef.Var) typedef.Var {
	castedArgs := make([]value.Value, 0)
	for i, arg := range args {
		if i >= len(f.Sig.Params) {
			castedArgs = append(castedArgs, arg.Load(bh))
			continue
		}

		expected := f.Sig.Params[i]
		raw := typeHandler.ImplicitTypeCast(bh, utils.GetTypeString(expected), arg.Load(bh))
		castedArgs = append(castedArgs, raw)
	}
	result := bh.N.NewCall(f, castedArgs...)
	return typeHandler.BuildVar(bh, typedef.NewType(utils.GetTypeString(result.Type())), result)
}
