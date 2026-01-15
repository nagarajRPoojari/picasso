package funcs

import (
	"fmt"

	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/identifier"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/state"
)

// DeclareFunc registers a method's signature within the LLVM module and the class metadata.
// It establishes the function's identity, including its name-mangled identifier,
// return type, and the formal parameter list—crucially appending the 'this'
// pointer to enable instance-aware execution.
//
// Technical Logic:
//   - Parameter Mapping: Converts AST parameter definitions into concrete LLVM IR
//     parameters by resolving types through the TypeHandler.
//   - This-Injection: Implements the Picasso ABI by adding a final parameter
//     representing the class UDT, allowing methods to access instance fields.
//   - Name Mangling: Uses the IdentifierBuilder to generate a unique, fully-qualified
//     name (e.g., "ClassName.MethodName") to avoid global symbol collisions.
//   - Memoization: Checks the existing Method map to prevent redundant declarations
//     and stores the resulting function symbol for the definition pass.
func (t *FuncHandler) DeclareFunc(cls string, st ast.FunctionDefinitionStatement, sourcePkg state.PackageEntry) {
	fqClsName := identifier.NewIdentifierBuilder(sourcePkg.Name).Attach(cls)
	aliasClsName := identifier.NewIdentifierBuilder(sourcePkg.Alias).Attach(cls)

	params := make([]*ir.Param, 0)
	argsTypes := make([]ast.Type, 0)
	for _, p := range st.Parameters {
		argsTypes = append(argsTypes, p.Type)
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
		}
		t.st.Classes[aliasClsName].Methods[aliasFuncName] = f
		t.st.Classes[aliasClsName].Returns[aliasFuncName] = st.ReturnType
	}

	if _, ok := t.st.Classes[aliasClsName].MethodArgs[aliasFuncName]; !ok {
		t.st.Classes[aliasClsName].MethodArgs[aliasFuncName] = argsTypes
	}
}
