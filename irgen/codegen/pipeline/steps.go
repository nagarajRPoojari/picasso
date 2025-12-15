package pipeline

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/libs"
	function "github.com/nagarajRPoojari/niyama/irgen/codegen/libs/func"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

func (t *Pipeline) importModules(methodMap map[string]function.Func, module string) {
	mod, ok := libs.ModuleList[module]
	if !ok {
		errorutils.Abort(errorutils.UnknownModule, module)
	}
	for name, f := range mod.ListAllFuncs() {
		n := fmt.Sprintf("%s.%s", module, name)
		methodMap[n] = f
	}
}

func (t *Pipeline) DeclareGlobals() {
	Loop(t.tree, func(st ast.VariableDeclarationStatement) {
		errorutils.Abort(errorutils.GlobalVarsNotAllowedError)
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
		class.ClassHandlerInst.DeclareClassUDT(st)
	})
}

func (t *Pipeline) Register() {
	tp := []string{c.TYPE_ATOMIC_BOOL, c.TYPE_ATOMIC_CHAR, c.TYPE_ATOMIC_SHORT, c.TYPE_ATOMIC_INT}
	for _, tpc := range tp {
		udt := t.st.CI.Types[tpc]
		t.st.Module.NewTypeDef(tpc, udt)
		mc := &typedef.MetaClass{
			FieldIndexMap:     make(map[string]int),
			ArrayVarsEleTypes: make(map[int]types.Type),
			VarAST:            make(map[string]*ast.VariableDeclarationStatement),
			UDT:               types.NewPointer(udt),
			Methods:           make(map[string]*ir.Func),
			Returns:           map[string]ast.Type{},
		}
		t.st.Classes[tpc] = mc
		t.st.TypeHandler.Register(tpc, mc)
	}
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
			f := t.st.Module.NewFunc(constants.MAIN, types.NewPointer(types.I8), ir.NewParam("", types.NewPointer(types.I8)))
			t.st.MainFunc = f
			funcs.FuncHandlerInst.DefineMainFunc(&st, make(map[string]struct{}))
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
