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
	classDefs := make(map[string]ast.ClassDeclarationStatement)
	roots := make([]ast.ClassDeclarationStatement, 0)

	for _, i := range t.tree.Body {
		if st, ok := i.(ast.ClassDeclarationStatement); ok {
			aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(st.Name)
			roots = append(roots, st)
			classDefs[aliasClsName] = st
		}
	}

	t.st.TypeHeirarchy = state.TypeHeirarchy{
		Roots:     roots,
		ClassDefs: classDefs,
	}

	for _, i := range roots {
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

func (t *Pipeline) declareFuncs(sourcePkg state.PackageEntry) {
	for _, i := range t.st.TypeHeirarchy.Roots {
		class.ClassHandlerInst.DeclareFunctions(i, sourcePkg)
	}
}

func (t *Pipeline) defineClasses() {
	for _, i := range t.st.TypeHeirarchy.Roots {
		class.ClassHandlerInst.DefineClass(i)
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
