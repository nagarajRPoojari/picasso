package funcs

import (
	"github.com/llir/llvm/ir/types"

	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

func (t *FuncHandler) DeclareFunc(cls string, st ast.FunctionDefinitionStatement) {
	params := make([]*ir.Param, 0)
	for _, p := range st.Parameters {
		params = append(params, ir.NewParam(p.Name, t.st.TypeHandler.GetLLVMType(tf.Type(p.Type.Get()))))
	}

	// at the end pass `this` parameter representing current object
	udt := t.st.Classes[cls].UDT
	params = append(params, ir.NewParam(constants.THIS, udt))

	name := t.st.IdentifierBuilder.Attach(cls, st.Name)

	var retType types.Type
	if st.ReturnType != nil {
		retType = t.st.TypeHandler.GetLLVMType(tf.Type(st.ReturnType.Get()))
	} else {
		retType = t.st.TypeHandler.GetLLVMType(tf.Type(tf.NULL))
	}

	if _, ok := t.st.Classes[cls].Methods[name]; !ok {
		f := t.st.Module.NewFunc(name, retType, params...)
		t.st.Classes[cls].Methods[name] = f
	}
}
