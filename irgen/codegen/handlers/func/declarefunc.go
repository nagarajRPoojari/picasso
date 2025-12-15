package funcs

import (
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
)

// DeclareFunc declares a class method in the IR module.
// It sets up the method's parameters and determines its return type.
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

	if _, ok := t.st.Classes[cls].Methods[name]; !ok {
		f := t.st.Module.NewFunc(name, retType, params...)
		t.st.Classes[cls].Methods[name] = f
		t.st.Classes[cls].Returns[name] = st.ReturnType
	}
}
