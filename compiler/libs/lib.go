package libs

import (
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs/io"
)

type Module interface {
	ListAllFuncs() map[string]function.Func
}

var ModuleList = make(map[string]Module)

func init() {
	ModuleList["io"] = io.NewIO()
}
