package libs

import (
	"github.com/nagarajRPoojari/x-lang/compiler/libs/array"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/io"
	types "github.com/nagarajRPoojari/x-lang/compiler/libs/type"
)

type Module interface {
	ListAllFuncs() map[string]function.Func
}

var ModuleList = make(map[string]Module)

func init() {
	ModuleList["io"] = io.NewIO()
	ModuleList["types"] = types.NewTypeHandler()
	ModuleList["array"] = array.NewArrayHandler()
}
