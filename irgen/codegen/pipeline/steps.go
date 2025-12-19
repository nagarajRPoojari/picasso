package pipeline

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

func (t *Pipeline) predeclareClasses(sourcePkg state.PackageEntry) {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		class.ClassHandlerInst.DeclareClassUDT(st, sourcePkg)
	})
}

func (t *Pipeline) registerTypes() {
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

func (t *Pipeline) declareVars(sourcePkg state.PackageEntry) {
	parent := make(map[string]string)
	classDefs := make(map[string]ast.ClassDeclarationStatement)
	childs := make(map[string][]ast.ClassDeclarationStatement)
	roots := make([]ast.ClassDeclarationStatement, 0)
	orphans := make([]ast.ClassDeclarationStatement, 0)

	for _, i := range t.tree.Body {
		if st, ok := i.(ast.ClassDeclarationStatement); ok {
			aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(st.Name)

			parent[aliasClsName] = st.Implements
			if st.Implements != "" {
				if _, ok := childs[st.Implements]; !ok {
					childs[st.Implements] = make([]ast.ClassDeclarationStatement, 0)
				}
				childs[st.Implements] = append(childs[st.Implements], st)
			} else {
				roots = append(roots, st)
			}
			classDefs[aliasClsName] = st
		}
	}

	// orphans are the one whose parent class resides in different module. so its
	// parent class ast is unavailble. I can safely traverse without bothering about its
	// order since it's parent class (from imported module) is already declared.
	for i := range classDefs {
		if _, ok := classDefs[parent[i]]; !ok {
			orphans = append(orphans, classDefs[i])
		}
	}

	// for inheritance involving only current packages, I can safely check its
	// cyclic inheritance condition. @todo: need to check cyclic inheritance
	// involving classes from other modules.
	for i := range parent {
		cyclicCheck(i, parent, make(map[string]struct{}))
	}

	t.st.TypeHeirarchy = state.TypeHeirarchy{
		Parent:    parent,
		Roots:     roots,
		Childs:    childs,
		ClassDefs: classDefs,
		Orphans:   orphans,
	}

	for _, i := range roots {
		traverse(i, childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DefineClassUDT(st, sourcePkg)
		})
	}
	for _, i := range orphans {
		class.ClassHandlerInst.DefineClassUDT(i, sourcePkg)
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

func (t *Pipeline) declareFuncs(sourcePkg state.PackageEntry) {
	for _, i := range t.st.TypeHeirarchy.Roots {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DeclareFunctions(st, sourcePkg)
		})
	}
	for _, i := range t.st.TypeHeirarchy.Orphans {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DeclareFunctions(st, sourcePkg)
		})
	}
}

func (t *Pipeline) defineClasses() {
	for _, i := range t.st.TypeHeirarchy.Roots {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DefineClass(st)
		})
	}
	for _, i := range t.st.TypeHeirarchy.Orphans {
		traverse(i, t.st.TypeHeirarchy.Childs, func(st ast.ClassDeclarationStatement) {
			class.ClassHandlerInst.DefineClass(st)
		})
	}
}

func (t *Pipeline) defineMain() {
	Loop(t.tree, func(st ast.FunctionDefinitionStatement) {
		if st.Name == constants.MAIN {
			f := t.st.Module.NewFunc(constants.MAIN, types.NewPointer(types.I8), ir.NewParam("", types.NewPointer(types.I8)))
			t.st.MainFunc = f
			funcs.FuncHandlerInst.DefineMainFunc(&st, make(map[string]struct{}))
		}
	})
}

func (t *Pipeline) insertYields() {
	m := t.st.Module
	// Define or get the yield function
	yieldFunc := c.Instance.Funcs[c.FUNC_SELF_YIELD]

	for _, fn := range m.Funcs {
		// Skip the yield function itself
		if fn.Name() == yieldFunc.Name() {
			continue
		}

		for _, blk := range fn.Blocks {
			var newInsts []ir.Instruction

			for _, inst := range blk.Insts {
				// If this is a call instruction, insert yield BEFORE it
				switch inst.(type) {
				case *ir.InstCall:
					newInsts = append(newInsts, ir.NewCall(yieldFunc))
				}

				// Then append the original instruction
				newInsts = append(newInsts, inst)
			}

			blk.Insts = newInsts
		}
	}
}

func Loop[T ast.Statement](tree ast.BlockStatement, fn func(T)) {
	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case T:
			fn(st)
		}
	}
}
