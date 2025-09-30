package pipeline

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/class"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/state"
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

func (t *Pipeline) DeclareInterface() {
	// Loop(t.tree, func(st ast.ClassDeclarationStatement) {
	// 	class.ClassHandlerInst.DeclareClassUDT(st)
	// })
}

func (t *Pipeline) PredeclareClasses() {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.DeclareClassUDT(st)
	})
}

func (t *Pipeline) DeclareVars() {
	parent := make(map[string]string)
	classDefs := make(map[string]ast.ClassDeclarationStatement)
	childs := make(map[string][]ast.ClassDeclarationStatement)
	roots := make([]ast.ClassDeclarationStatement, 0)

	for _, i := range t.tree.Body {
		if st, ok := i.(ast.ClassDeclarationStatement); ok {
			parent[st.Name] = st.Implements
			if st.Implements != "" {
				if _, ok := childs[st.Implements]; !ok {
					childs[st.Implements] = make([]ast.ClassDeclarationStatement, 0)
				}
				childs[st.Implements] = append(childs[st.Implements], st)
			} else {
				roots = append(roots, st)
			}
			classDefs[st.Name] = st
		}
	}
	for i := range parent {
		cyclicCheck(i, parent, make(map[string]struct{}))
	}

	t.st.TypeHeirarchy = state.TypeHeirarchy{
		Parent:    parent,
		Roots:     roots,
		Childs:    childs,
		ClassDefs: classDefs,
	}

	for _, i := range roots {
		traverse(i, childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DefineClassUDT(st)
		})
	}

}

func cyclicCheck(child string, parent map[string]string, isV map[string]struct{}) {
	isV[child] = struct{}{}
	p := parent[child]
	if p != "" {
		if _, ok := isV[p]; ok {
			panic(fmt.Sprintf("cyclic inheritance found involving %s", p))
		}
		cyclicCheck(p, parent, isV)
	}
	delete(isV, child)
}

func traverse(parent ast.ClassDeclarationStatement, childs map[string][]ast.ClassDeclarationStatement, fn func(ast.ClassDeclarationStatement)) {
	fn(parent)
	for _, c := range childs[parent.Name] {
		traverse(c, childs, fn)
	}
}

func (t *Pipeline) DeclareFuncs() {
	for _, i := range t.st.TypeHeirarchy.Roots {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DeclareFunctions(st)
		})
	}
}

func (t *Pipeline) DefineClasses() {
	for _, i := range t.st.TypeHeirarchy.Roots {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DefineClass(st)
		})
	}
}

func (t *Pipeline) DefineMain() {
	Loop(t.tree, func(st ast.FunctionDefinitionStatement) {
		if st.Name == constants.MAIN {
			f := t.st.Module.NewFunc(constants.MAIN, types.I32)
			t.st.MainFunc = f
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
