package class

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

// DeclareClassUDT registers a new User-Defined Type (UDT) within the LLVM module.
// It initializes the class as an opaque struct to allow for recursive type
// definitions (pointers to self) and prepares a MetaClass container to hold
// field offsets, method symbols, and type metadata for semantic resolution.
//
// Key Logic:
//   - Validates class uniqueness to prevent symbol collisions.
//   - Defines a named opaque struct in the LLVM module.
//   - Maps the class name to its MetaClass metadata in the global state for
//     future field lookups and method dispatch.
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

// DeclareFunctions orchestrates the declaration of all member functions
// associated with a class. It performs a pass over the class body and
// inherited definitions to populate the function symbol table before
// actual function bodies are emitted.
//
// Key Logic:
//   - Iterates through the local AST body to register member functions.
//   - Traverses the Type Hierarchy to pull in parent class method signatures,
//     facilitating inheritance and polymorphism in the generated IR.
//   - Delegates signature creation to the FuncHandler to ensure consistent
//     ABI naming and parameter lowering.
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
