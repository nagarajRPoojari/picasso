package class

import (
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
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
func (t *ClassHandler) DeclareClassUDT(cls ast.ClassDeclarationStatement, sourcePkg state.PackageEntry) {

	clsName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(cls.Name)
	aliasName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(cls.Name)

	if _, ok := t.st.Classes[aliasName]; ok {
		errorutils.Abort(errorutils.TypeRedeclaration, clsName)
	}

	// define as opaque type
	udt := types.NewStruct()
	if _, ok := t.st.GlobalTypeList[clsName]; !ok {
		t.st.GlobalTypeList[clsName] = t.st.Module.NewTypeDef(clsName, udt)
	}
	mc := tf.NewMetaClass(types.NewPointer(udt), cls.Implements)
	t.st.Classes[aliasName] = mc

	// assumed that all interfaces are defined first
	if cls.Implements != "" {
		errorutils.Assert(t.st.Interfaces != nil, "interfaces USTs are expected to be initialized before class UDT")
		t.st.Interfaces[cls.Implements].ImplementedBy = append(t.st.Interfaces[cls.Implements].ImplementedBy, aliasName)
	}

	// register current class type with TypeHandler. this allows current class
	// to be identified as a valid type in future while building vars & type
	// conversions.
	t.st.TypeHandler.RegisterClass(aliasName, mc)

	t.st.Classes[aliasName].Internal = cls.IsInternal
}

// DeclareFunctions orchestrates the declaration of all member functions
// associated with a class. It performs a pass over the class body and
// inherited definitions to populate the function symbol table before
// actual function bodies are emitted.
//
// Key Logic:
//   - Iterates through the local AST body to register member functions.
//   - Delegates signature creation to the FuncHandler to ensure consistent
//     ABI naming and parameter lowering.
func (t *ClassHandler) DeclareFunctions(cls ast.ClassDeclarationStatement, sourcePkg state.PackageEntry) {
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DeclareFunc(cls.Name, st, sourcePkg)
		}
	}
}
