package pipeline

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/class"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
	"github.com/nagarajRPoojari/x-lang/compiler/libs"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

func (t *Pipeline) importModules(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("unable to find module: %s", module))
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s", module, name)
		methodMap[n] = f
	}
}

func (t *Pipeline) DeclareGlobals() {
	Loop(t.tree, func(st ast.VariableDeclarationStatement) {
		errorsx.PanicCompilationError("global vars not allowed")
	})
}

func (t *Pipeline) ImportModules() {
	Loop(t.tree, func(st ast.ImportStatement) {
		if st.From == constants.BUILTIN {
			t.importModules(t.st.LibMethods, st.Name)
		}
	})
}

func (t *Pipeline) PredeclareClasses() {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.PredeclareClass(st)
	})
}

func (t *Pipeline) DeclareVars() {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.DefineClassVars(st)
	})
}

func (t *Pipeline) DeclareFuncs() {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.DeclareFunctions(st)
	})
}

func (t *Pipeline) DefineClasses() {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.DefineClass(st)
	})
}

func (t *Pipeline) DefineMain() {
	Loop(t.tree, func(st ast.FunctionDeclarationStatement) {
		if st.Name == constants.MAIN {
			f := t.st.Module.NewFunc(constants.MAIN, types.I32)
			t.st.Methods[constants.MAIN] = f
			funcs.FuncHandlerInst.DefineFunc("", &st)
		}
	})
}

func Loop[T ast.Statement](tree ast.BlockStatement, fn func(T)) {
	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case T:
			fn(st)
		}
	}
}
