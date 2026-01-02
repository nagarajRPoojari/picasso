package class

import (
	"fmt"

	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	funcs "github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/func"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	typedef "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
)

// DefineClass triggers the emission of concrete LLVM IR function bodies for
// all methods belonging to a class. It handles the traversal of local
// definitions using a tracking map to ensure that overridden methods are
// defined only once using the most specific implementation.
func (t *ClassHandler) DefineClass(cls ast.ClassDeclarationStatement) {
	avoid := make(map[string]struct{}, 0)
	fqClsName := t.st.IdentifierBuilder.Attach(cls.Name)
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			funcs.FuncHandlerInst.DefineFunc(fqClsName, &st, avoid)
			avoid[st.Name] = struct{}{}
		}
	}
}

// DefineClassUDT resolves a previously declared opaque struct into a concrete
// LLVM struct layout. It maps class members (fields and methods) to numerical
// indices for GEP (GetElementPtr) instructions and populates the struct's
// memory footprint.
//
// Key Logic:
//   - Interface Validation: If the class implements an interface, it verifies
//     that all required methods exist and their signatures (hashed) match
//     the interface definition.
//   - Member Mapping: Iterates through the class body to assign indices to
//     variables and function pointers, storing metadata in the ClassMeta (mc).
//   - Type Resolution: Converts AST types into concrete LLVM types for each
//     struct field.
//   - Opaque Completion: Updates the underlying LLVM struct type (stored in
//     mc.UDT) with the finalized field list, completing the type definition.
func (t *ClassHandler) DefineClassUDT(cls ast.ClassDeclarationStatement, sourcePkg state.PackageEntry) {
	// fqClsName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(cls.Name)
	aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(cls.Name)

	mc := t.st.Classes[aliasClsName]

	fieldTypes := make([]types.Type, 0)
	vars := make(map[string]struct{}, 0)

	// map each fields with corresponding udt struct index
	i := 0

	shouldImplement := map[string]typedef.MethodSig{}
	if cls.Implements != "" {
		interfaceName := cls.Implements
		if _, ok := t.st.Interfaces[interfaceName]; !ok {
			errorutils.Abort(errorutils.UnknownInterfaceError, interfaceName)
		}

		ifMeta := t.st.Interfaces[interfaceName]
		shouldImplement = ifMeta.Methods
	}

	// Opaque resolution: define concrete types of all fields
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			aliasVarName := fmt.Sprintf("%s.%s", aliasClsName, st.Identifier)
			if _, ok := vars[aliasVarName]; ok {
				errorutils.Abort(errorutils.VariableRedeclaration, st.Identifier)
			}

			mc.FieldIndexMap[aliasVarName] = i
			mc.VarAST[aliasVarName] = &st

			fieldType := t.st.TypeHandler.GetLLVMType(st.ExplicitType.Get())
			fieldTypes = append(fieldTypes, fieldType)
			vars[aliasVarName] = struct{}{}

			if st.ExplicitType.GetUnderlyingType() != "" {
				// mc.ArrayVarsEleTypes is a map of field index to its underlying element type
				mc.ArrayVarsEleTypes[i] = t.st.TypeHandler.GetLLVMType(st.ExplicitType.GetUnderlyingType())
			}

			if st.IsInternal {
				mc.InternalFields[aliasVarName] = struct{}{}
			}

			i++

		case ast.FunctionDefinitionStatement:
			aliasFuncName := fmt.Sprintf("%s.%s", aliasClsName, st.Name)
			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
			} else {
				// empty string is expected to give a void type
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			args := make([]types.Type, 0)
			for _, p := range st.Parameters {
				args = append(args, t.st.TypeHandler.GetLLVMType(p.Type.Get()))
			}

			if sh, ok := shouldImplement[st.Name]; ok {
				if sh.Hash != utils.HashFuncSig(st.Parameters, st.ReturnType) {
					errorutils.Abort(errorutils.FunctionSignatureMisMatch, aliasFuncName)
				}
				delete(shouldImplement, st.Name)
			}

			// args = append(args, mc.UDT)
			funcType := types.NewFunc(retType, args...)
			fieldTypes = append(fieldTypes, types.NewPointer(funcType))

			mc.FieldIndexMap[aliasFuncName] = i

			if st.IsInternal {
				mc.InternalFields[aliasFuncName] = struct{}{}
			}

			i++
		}
	}

	if len(shouldImplement) > 0 {
		errorutils.Abort(errorutils.UnImplementedInterfaceMethod, shouldImplement)
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
