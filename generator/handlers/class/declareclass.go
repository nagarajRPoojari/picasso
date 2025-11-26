package class

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	funcs "github.com/nagarajRPoojari/x-lang/generator/handlers/func"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
)

// DeclareClassUDT registers a new user-defined class type in the module.
// It creates an opaque LLVM struct for the class, wraps it in a pointer type,
// and stores metadata including fields, variables, and methods.
//
// If the class name already exists in the symbol table, it does nothing.
//
// Params:
//
//	cls – the AST ClassDeclarationStatement defining the class
func (t *ClassHandler) DeclareClassUDT(cls ast.ClassDeclarationStatement) {
	if _, ok := t.st.Classes[cls.Name]; ok {
		errorutils.Abort(errorutils.ClassRedeclaration, cls.Name)
	}
	udt := types.NewStruct() // opaque
	t.st.Module.NewTypeDef(cls.Name, udt)
	mc := &tf.MetaClass{
		FieldIndexMap:     make(map[string]int),
		ArrayVarsEleTypes: make(map[int]types.Type),
		VarAST:            make(map[string]*ast.VariableDeclarationStatement),
		UDT:               types.NewPointer(udt),
		Methods:           make(map[string]*ir.Func),
		Returns:           map[string]ast.Type{},
	}
	t.st.Classes[cls.Name] = mc
	t.st.TypeHandler.Register(cls.Name, mc)
}

// DeclareFunctions declares all functions (methods) of a class in the IR.
// It processes both the functions defined directly in the class body and
// any functions inherited from its parent class.
//
// Params:
//
//	cls – the AST ClassDeclarationStatement for which functions are declared
func (t *ClassHandler) DeclareFunctions(cls ast.ClassDeclarationStatement) {

	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DeclareFunc(cls.Name, st)
		}
	}

	for _, stI := range t.st.TypeHeirarchy.ClassDefs[t.st.TypeHeirarchy.Parent[cls.Name]].Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DeclareFunc(cls.Name, st)
		}
	}
}
