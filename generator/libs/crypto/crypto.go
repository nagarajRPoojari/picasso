package crypto

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/utils"

	function "github.com/nagarajRPoojari/x-lang/generator/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

type Crypto struct {
}

func NewCrypto() *Crypto {
	return &Crypto{}

}

func (t *Crypto) ListAllFuncs() map[string]function.Func {
	funcs := make(map[string]function.Func)
	funcs[c.HASH] = t.hash
	return funcs
}

func (t *Crypto) hash(typeHandler *tf.TypeHandler, module *ir.Module, bh *bc.BlockHolder, args []tf.Var) tf.Var {
	method := c.NewInterface(module).Funcs[c.HASH]
	castedArgs := make([]value.Value, 0)
	for i, arg := range args {
		if i >= len(method.Sig.Params) {
			castedArgs = append(castedArgs, arg.Load(bh))
			continue
		}
		castedArgs = append(castedArgs, arg.Load(bh))
	}
	result := bh.N.NewCall(method, castedArgs...)
	return typeHandler.BuildVar(bh, tf.NewType(utils.GetTypeString(result.Type())), result)
}
