package interfaceh

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
)

// DefineInterfaceUDT finalizes the interface's structural definition by resolving
// its opaque type into a concrete LLVM struct layout. This layout serves as
// a virtual method table (vtable) template for all classes that implement
// the interface.
//
// Key Logic:
//   - Structural Constraint Enforcement: Explicitly prevents the definition of
//     stateful fields (variables), as interfaces are restricted to method
//     signatures.
//   - Function Pointer Mapping: Converts method definitions into LLVM pointer-to-function
//     types, populating the struct's field list to facilitate dynamic dispatch.
//   - Field Indexing: Records the precise offset of each method within the struct
//     to enable GEP (GetElementPtr) lookup during interface calls.
//   - Opaque Resolution: Finalizes the previously declared opaque struct by
//     assigning the calculated fieldTypes, effectively closing the type definition.
func (t *InterfaceHandler) DefineInterfaceUDT(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {
	fqInterfaceName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

	mc := t.st.Classes[fqInterfaceName]

	fieldTypes := make([]types.Type, 0)
	i := 0

	// Opaque resolution: define concrete types of all fields
	for _, stI := range ifs.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			errorutils.Abort(errorutils.VarsNotAllowedInInterfaceError)

		case ast.FunctionDefinitionStatement:
			fqFuncName := fmt.Sprintf("%s.%s", fqInterfaceName, st.Name)
			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(t.st.ResolveAlias(st.ReturnType.Get()))
			} else {
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			args := make([]types.Type, 0)
			for _, p := range st.Parameters {
				args = append(args, t.st.TypeHandler.GetLLVMType(t.st.ResolveAlias(p.Type.Get())))
			}

			funcType := types.NewFunc(retType, args...)
			fieldTypes = append(fieldTypes, types.NewPointer(funcType))

			mc.FieldIndexMap[fqFuncName] = i
			i++
		}
	}

	ptr, ok := mc.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be a pointer")
	}
	st, ok := ptr.ElemType.(*types.StructType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalUDTDefinitionError, "udt must be pointer to a struct")
	}
	st.Fields = fieldTypes
}
