package libs

import (
	"github.com/nagarajRPoojari/x-lang/generator/libs/array"
	sync "github.com/nagarajRPoojari/x-lang/generator/libs/atomic"
	"github.com/nagarajRPoojari/x-lang/generator/libs/crypto"
	function "github.com/nagarajRPoojari/x-lang/generator/libs/func"
	"github.com/nagarajRPoojari/x-lang/generator/libs/io"
	"github.com/nagarajRPoojari/x-lang/generator/libs/strings"
	types "github.com/nagarajRPoojari/x-lang/generator/libs/type"
)

type Module interface {
	ListAllFuncs() map[string]function.Func
}

var ModuleList = make(map[string]Module)

func init() {
	ModuleList["syncio"] = io.NewSyncIO()
	ModuleList["crypto"] = crypto.NewCrypto()
	ModuleList["asyncio"] = io.NewAsyncIO()
	ModuleList["types"] = types.NewTypeHandler()
	ModuleList["array"] = array.NewArrayHandler()
	ModuleList["sync"] = sync.NewSync()
	ModuleList["strings"] = strings.NewStringHandler()
}
