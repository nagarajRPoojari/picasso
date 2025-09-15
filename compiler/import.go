package compiler

import (
	"fmt"

	"github.com/nagarajRPoojari/x-lang/compiler/libs"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *LLVM) importer(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("unable to find module: %s", module))
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s", module, name)
		methodMap[n] = f
	}
}
