package libs

import (
	"github.com/nagarajRPoojari/picasso/irgen/codegen/libs/array"
	function "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/func"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/libs/io"
	types "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/type"
)

type Module interface {
	ListAllFuncs() map[string]function.Func
}

var ModuleList = make(map[string]Module)

func init() {
	ModuleList["types"] = types.NewTypeHandler()
	ModuleList["array"] = array.NewArrayHandler()
	ModuleList["syncio"] = io.NewSyncIO()
}
