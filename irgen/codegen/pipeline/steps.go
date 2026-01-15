package pipeline

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/c"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/class"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	funcs "github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/func"
	interfaceh "github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/interface"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
	typedef "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
)

func (t *Pipeline) predeclareClasses(sourcePkg state.PackageEntry) {
	Loop(t.tree, func(st ast.ClassDeclarationStatement) {
		t.m.GetClassHandler().(*class.ClassHandler).DeclareOpaqueClass(st, sourcePkg)
	})
}

func (t *Pipeline) predeclareInterfraces(sourcePkg state.PackageEntry) {
	Loop(t.tree, func(st ast.InterfaceDeclarationStatement) {
		t.m.GetInterfaceHandler().(*interfaceh.InterfaceHandler).DeclareInterface(st, sourcePkg)
	})
}

func (t *Pipeline) registerTypes() {
	for tpc, udt := range t.st.CI.Types {
		t.st.Module.NewTypeDef(tpc, udt)
		mc := &typedef.MetaClass{
			FieldIndexMap:     make(map[string]int),
			ArrayVarsEleTypes: make(map[int]types.Type),
			InternalFields:    make(map[string]struct{}),
			VarAST:            make(map[string]*ast.VariableDeclarationStatement),
			UDT:               types.NewPointer(udt),
			Methods:           make(map[string]*ir.Func),
			MethodArgs:        make(map[string][]ast.Type),
			Returns:           map[string]ast.Type{},
		}
		t.st.Classes[tpc] = mc
		t.st.TypeHandler.RegisterClass(tpc, mc)
	}
}

func (t *Pipeline) declareClassFields(sourcePkg state.PackageEntry) {
	roots := make([]ast.ClassDeclarationStatement, 0)

	for _, i := range t.tree.Body {
		if st, ok := i.(ast.ClassDeclarationStatement); ok {
			roots = append(roots, st)
		}
	}

	t.st.TypeHeirarchy.ClassRoots = roots
	for _, i := range roots {
		t.m.GetClassHandler().(*class.ClassHandler).DefineClass(i, sourcePkg)
	}
}

func (t *Pipeline) declareInterfaceFields(sourcePkg state.PackageEntry) {
	roots := make([]ast.InterfaceDeclarationStatement, 0)

	for _, i := range t.tree.Body {
		if st, ok := i.(ast.InterfaceDeclarationStatement); ok {
			roots = append(roots, st)
		}
	}

	t.st.TypeHeirarchy.InterfaceRoots = roots
	for _, i := range roots {
		t.m.GetInterfaceHandler().(*interfaceh.InterfaceHandler).DefineInterfaceUDT(i, sourcePkg)
	}
}

func (t *Pipeline) declareClassFuncs(sourcePkg state.PackageEntry) {
	for _, i := range t.st.TypeHeirarchy.ClassRoots {
		t.m.GetClassHandler().(*class.ClassHandler).DeclareClassFuncs(i, sourcePkg)
	}
}

func (t *Pipeline) declareInterfaceFuncs(sourcePkg state.PackageEntry) {
	for _, i := range t.st.TypeHeirarchy.InterfaceRoots {
		t.m.GetInterfaceHandler().(*interfaceh.InterfaceHandler).DeclareClassFuncs(i, sourcePkg)
	}
}

func (t *Pipeline) defineClasses() {
	for _, i := range t.st.TypeHeirarchy.ClassRoots {
		t.m.GetClassHandler().(*class.ClassHandler).DefineClassFuncs(i)
	}
}

func (t *Pipeline) defineMain() {
	Loop(t.tree, func(st ast.FunctionDefinitionStatement) {
		if st.Name == constants.MAIN {
			f := t.st.Module.NewFunc(constants.MAIN, types.NewPointer(types.I8), ir.NewParam("", types.NewPointer(types.I8)))
			t.st.MainFunc = f
			t.m.GetFuncHandler().(*funcs.FuncHandler).DefineMainFunc(&st, make(map[string]struct{}))
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
