package funcs

import (
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
)

// DeclareFunc registers a method's signature within the LLVM module and the class metadata.
// It establishes the function's identity, including its name-mangled identifier,
// return type, and the formal parameter list—crucially appending the 'this'
// pointer to enable instance-aware execution.
//
// Technical Logic:
//   - Parameter Mapping: Converts AST parameter definitions into concrete LLVM IR
//     parameters by resolving types through the TypeHandler.
//   - This-Injection: Implements the Niyama ABI by adding a final parameter
//     representing the class UDT, allowing methods to access instance fields.
//   - Name Mangling: Uses the IdentifierBuilder to generate a unique, fully-qualified
//     name (e.g., "ClassName.MethodName") to avoid global symbol collisions.
//   - Memoization: Checks the existing Method map to prevent redundant declarations
//     and stores the resulting function symbol for the definition pass.
func (t *FuncHandler) DeclareFunc(cls string, st ast.FunctionDefinitionStatement) {
	params := make([]*ir.Param, 0)
	for _, p := range st.Parameters {
		params = append(params, ir.NewParam(p.Name, t.st.TypeHandler.GetLLVMType(p.Type.Get())))
	}

	// at the end pass `this` parameter representing current object
	udt := t.st.Classes[cls].UDT
	params = append(params, ir.NewParam(constants.THIS, udt))

	name := t.st.IdentifierBuilder.Attach(cls, st.Name)

	var retType types.Type
	if st.ReturnType != nil {
		retType = t.st.TypeHandler.GetLLVMType(st.ReturnType.Get())
	} else {
		retType = t.st.TypeHandler.GetLLVMType("")
	}

	// store current functions so that later during class instantiation instance
	// can be made pointing to the functions.
	if _, ok := t.st.Classes[cls].Methods[name]; !ok {
		f := t.st.Module.NewFunc(name, retType, params...)
		t.st.Classes[cls].Methods[name] = f
		t.st.Classes[cls].Returns[name] = st.ReturnType
	}
}
