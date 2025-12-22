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
// all methods belonging to a class. It handles the traversal of both local
// definitions and inherited parent methods, using a tracking map to ensure
// that overridden methods are defined only once using the most specific
// implementation.
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

// DefineClassUDT performs the "lowering" of high-level class definitions into
// concrete LLVM struct layouts. It handles the structural aspects of inheritance
// by flattening parent fields into the child struct, mapping field identifiers
// to numerical indices (for GEP instructions), and calculating the memory
// footprint for function pointers used in method dispatch.
//
// Key Logic:
//   - Structural Inheritance: Deeply copies field types from the parent UDT to
//     ensure binary compatibility for polymorphism.
//   - Member Indexing: Assigns monotonically increasing indices to fields and
//     methods to facilitate GetElementPtr (GEP) offset calculations.
//   - Signature Validation: Uses method hashes to verify that overrides match
//     the parent signature, aborting on interface mismatches.
//   - Opaque Resolution: Updates the previously declared opaque struct with
//     the finalized field set, completing the type definition in the LLVM module.
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

	fmt.Printf("cls.Implements: %v\n", cls.Implements)
	fmt.Printf("shouldImplement: %v\n", shouldImplement)

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
				mc.ArrayVarsEleTypes[i] = t.st.TypeHandler.GetLLVMType(st.ExplicitType.GetUnderlyingType())
			}
			i++

		case ast.FunctionDefinitionStatement:
			aliasFuncName := fmt.Sprintf("%s.%s", aliasClsName, st.Name)
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
