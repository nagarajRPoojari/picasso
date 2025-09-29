package class

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// defineClass similar to declareClass but does function concrete declaration
func (t *ClassHandler) DefineClassFuncs(cls ast.ClassDeclarationStatement) {
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DefineFunc(cls.Name, &st)
		}
	}
}

// defineClassVars stores corresponding ast for all var declaration
// which will be used to instantiate them on constructor call, i.e, new MyClass()
func (t *ClassHandler) DefineClassVars(cls ast.ClassDeclarationStatement) {
	mc := t.st.Classes[cls.Name]
	fieldTypes := make([]types.Type, 0)
	vars := make(map[string]struct{}, 0)

	// map each fields with corresponding udt struct index
	i := 0
	for _, stI := range cls.Body {
		if st, ok := stI.(ast.VariableDeclarationStatement); ok {
			fqName := t.st.IdentifierBuilder.Attach(cls.Name, st.Identifier)
			if _, ok := vars[fqName]; ok {
				errorsx.PanicCompilationError("variable already exists")
			}

			mc.VarIndexMap[fqName] = i
			mc.VarAST[fqName] = &st

			fieldType := t.st.TypeHandler.GetLLVMType(tf.Type(st.ExplicitType.Get()))
			fieldTypes = append(fieldTypes, fieldType)
			vars[fqName] = struct{}{}
			i++
		}
	}

	if ptr, ok := mc.UDT.(*types.PointerType); ok {
		if st, ok2 := ptr.ElemType.(*types.StructType); ok2 {
			st.Fields = fieldTypes
		} else {
			panic("mc.UDT pointer does not point to a struct")
		}
	} else {
		panic("mc.UDT is not a pointer type")
	}

}
