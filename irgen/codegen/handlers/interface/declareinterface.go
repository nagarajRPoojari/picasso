package interfaceh

import (
	"fmt"

	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/utils"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
)

// DeclareInterface registers a symbolic interface type within the global state
// and the LLVM module. Like classes, interfaces are initially declared as
// opaque structs to allow for mutually recursive type definitions and
// cross-package references.
//
// Key Logic:
//   - Namespace Validation: Ensures the interface name does not collide with
//     existing classes or interfaces in the current package scope.
//   - Opaque Type Allocation: Registers a named struct in LLVM to serve as the
//     base type for interface-to-concrete-type casting.
//   - Metadata Initialization: Populates a MetaInterface container to track
//     method signatures and the list of implementing classes for late-bound
//     validation.
//   - Type Registration: Informs the TypeHandler of the new interface type,
//     enabling its use as a valid type for parameters and variables.
func (t *InterfaceHandler) DeclareInterface(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {

	ifName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(ifs.Name)
	aliasName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

	if _, ok := t.st.Classes[aliasName]; ok {
		errorutils.Abort(errorutils.TypeRedeclaration, ifName)
	}

	udt := types.NewStruct() // opaque
	if _, ok := t.st.GlobalTypeList[ifName]; !ok {
		t.st.GlobalTypeList[ifName] = t.st.Module.NewTypeDef(ifName, udt)
	}

	// interface is treated just like a class but instantiation is prevented
	mc := tf.NewMetaClass(types.NewPointer(udt), "")
	t.st.Classes[aliasName] = mc

	mi := tf.NewMetaInterface()
	mi.UDT = types.NewPointer(udt)

	t.st.Interfaces[aliasName] = mi

	// register current interface type with TypeHandler. this allows current interface
	// to be identified as a valid type in future while building vars & type
	// conversions.
	t.st.TypeHandler.RegisterInterface(aliasName, mi)
}

// DeclareClassFuncs populates the interface with its defined method signatures.
// It generates LLVM function prototypes for each method, including the
// injection of the implicit 'this' parameter to support polymorphism.
//
// Key Logic:
//   - Signature Registration: Converts AST function definitions into LLVM
//     function types, mapping return types and parameters to their IR equivalents.
//   - Implicit 'this' Injection: Appends a pointer to the interface's UDT as the
//     final parameter, facilitating method dispatch on the implementing instance.
//   - Method Hashing: Calculates and stores a unique hash of the function
//     signature (parameters and return type) to ensure strict type safety
//     when classes implement this interface.
//   - Global Symbol Management: Registers methods in the GlobalFuncList to
//     prevent duplicate linkage symbols and enable cross-module accessibility
func (t *InterfaceHandler) DeclareClassFuncs(ifs ast.InterfaceDeclarationStatement, sourcePkg state.PackageEntry) {
	aliasName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)
	for _, stI := range ifs.Body {
		switch st := stI.(type) {
		case ast.FunctionDefinitionStatement:
			// store function signature for interface

			fqClsName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(ifs.Name)
			aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(ifs.Name)

			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.st.TypeHandler.GetLLVMType(p.Type.Get())))
			}

			// at the end pass `this` parameter representing current object
			udt := t.st.Classes[aliasClsName].UDT
			params = append(params, ir.NewParam(constants.THIS, udt))

			fqFuncName := fmt.Sprintf("%s.%s", fqClsName, st.Name)
			aliasFuncName := fmt.Sprintf("%s.%s", aliasClsName, st.Name)

			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
			} else {
				retType = t.st.TypeHandler.GetLLVMType("")
			}

			// store current functions so that later during class instantiation instance
			// can be made pointing to the functions.
			if _, ok := t.st.Classes[aliasClsName].Methods[aliasFuncName]; !ok {
				f, ok := t.st.GlobalFuncList[fqFuncName]
				if !ok {
					f = t.st.Module.NewFunc(fqFuncName, retType, params...)
					t.st.GlobalFuncList[fqFuncName] = f

					// store method signature to validate implementation in
					// implementor classes
					t.st.Interfaces[aliasName].Methods[st.Name] = tf.MethodSig{
						Hash:     utils.HashFuncSig(st.Parameters, st.ReturnType),
						Name:     st.Name,
						FuncType: f,
					}
				}
				t.st.Classes[aliasClsName].Methods[aliasFuncName] = f
				t.st.Classes[aliasClsName].Returns[aliasFuncName] = st.ReturnType
			}
		}
	}
}
