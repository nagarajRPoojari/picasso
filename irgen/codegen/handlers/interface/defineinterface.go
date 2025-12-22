package interfaceh

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
)

func (t *InterfaceHandler) DefineInterfaceUDT(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {
	aliasInterfaceName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

	mc := t.st.Classes[aliasInterfaceName]

	fieldTypes := make([]types.Type, 0)
	i := 0

	// Opaque resolution: define concrete types of all fields
	for _, stI := range ifs.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			errorutils.Abort(errorutils.VarsNotAllowedInInterfaceError)

		case ast.FunctionDefinitionStatement:
			aliasFuncName := fmt.Sprintf("%s.%s", aliasInterfaceName, st.Name)
			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
			} else {
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			args := make([]types.Type, 0)
			for _, p := range st.Parameters {
				args = append(args, t.st.TypeHandler.GetLLVMType(p.Type.Get()))
			}

			funcType := types.NewFunc(retType, args...)
			fieldTypes = append(fieldTypes, types.NewPointer(funcType))

			mc.FieldIndexMap[aliasFuncName] = i
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
