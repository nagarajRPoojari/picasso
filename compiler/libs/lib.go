package libs

import (
	"github.com/nagarajRPoojari/x-lang/compiler/libs/array"
	sync "github.com/nagarajRPoojari/x-lang/compiler/libs/atomic"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/io"
	types "github.com/nagarajRPoojari/x-lang/compiler/libs/type"
)

type Module interface {
	ListAllFuncs() map[string]function.Func
}

var ModuleList = make(map[string]Module)

func init() {
	ModuleList["syncio"] = io.NewSyncIO()
	ModuleList["asyncio"] = io.NewAsyncIO()
	ModuleList["types"] = types.NewTypeHandler()
	ModuleList["array"] = array.NewArrayHandler()
	ModuleList["sync"] = sync.NewSync()
}
