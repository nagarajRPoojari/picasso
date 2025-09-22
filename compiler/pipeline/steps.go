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
	for _, stI := range t.tree.Body {
		switch stI.(type) {
		case ast.VariableDeclarationStatement:
			errorsx.PanicCompilationError("global vars not allowed")
		}
	}
}

func (t *Pipeline) ImportModules() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.ImportStatement:
			t.importModules(t.st.LibMethods, st.Name)
		}
	}
}

func (t *Pipeline) PredeclareClasses() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			class.ClassHandlerInst.PredeclareClass(st)
		}
	}
}

func (t *Pipeline) DeclareVars() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			class.ClassHandlerInst.DefineClassVars(st)
		}
	}
}

func (t *Pipeline) DeclareFuncs() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			class.ClassHandlerInst.DeclareFunctions(st)
		}
	}
}

func (t *Pipeline) DefineClasses() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			class.ClassHandlerInst.DefineClass(st)
		}
	}
}

func (t *Pipeline) DefineMain() {
	for _, stI := range t.tree.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			if st.Name == constants.MAIN {
				f := t.st.Module.NewFunc(constants.MAIN, types.I32)
				t.st.Methods[constants.MAIN] = f
				funcs.FuncHandlerInst.DefineFunc("", &st)
			}
		}
	}
}
