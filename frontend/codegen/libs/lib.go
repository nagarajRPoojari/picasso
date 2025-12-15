package libs

import (
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/array"
	sync "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/atomic"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/crypto"
	function "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/func"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/io"
	"github.com/nagarajRPoojari/niyama/frontend/codegen/libs/strings"
	types "github.com/nagarajRPoojari/niyama/frontend/codegen/libs/type"
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
