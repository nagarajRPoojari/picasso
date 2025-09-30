package class

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// predeclareClass creates an opaque struct for all classes defined by user
// and registers it with typehandler for identfying forward declaration
func (t *ClassHandler) DeclareClassUDT(cls ast.ClassDeclarationStatement) {
	if _, ok := t.st.Classes[cls.Name]; ok {
		return
	}
	udt := types.NewStruct() // opaque
	t.st.Module.NewTypeDef(cls.Name, udt)
	mc := &tf.MetaClass{
		FieldIndexMap: make(map[string]int),
		VarAST:        make(map[string]*ast.VariableDeclarationStatement),
		UDT:           types.NewPointer(udt),
		Methods:       make(map[string]*ir.Func),
	}
	t.st.Classes[cls.Name] = mc
	t.st.TypeHandler.Register(cls.Name, mc)

}

// declareFunctions loops over all functions inside Class & creates
// a header declaration
func (t *ClassHandler) DeclareFunctions(cls ast.ClassDeclarationStatement) {

	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DeclareFunc(cls.Name, st)
		}
	}

	// declare all inherited methods
	for _, stI := range t.st.TypeHeirarchy.ClassDefs[t.st.TypeHeirarchy.Parent[cls.Name]].Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DeclareFunc(cls.Name, st)
		}
	}
}
