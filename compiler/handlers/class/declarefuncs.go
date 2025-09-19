package class

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// declareFunctions loops over all functions inside Class & creates
// a header declaration
func (t *ClassHandler) DeclareFunctions(cls ast.ClassDeclarationStatement) {
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.st.TypeHandler.GetLLVMType(tf.Type(p.Type.Get()))))
			}

			// at the end pass `this` parameter representing current object
			udt := t.st.Classes[cls.Name].UDT
			params = append(params, ir.NewParam(constants.THIS, udt))

			name := t.st.IdentifierBuilder.Attach(cls.Name, st.Name)

			var retType types.Type
			if st.ReturnType != nil {
				retType = t.st.TypeHandler.GetLLVMType(tf.Type(st.ReturnType.Get()))
			} else {
				retType = t.st.TypeHandler.GetLLVMType(tf.Type(tf.NULL))
			}
			f := t.st.Module.NewFunc(name, retType, params...)
			t.st.Methods[name] = f
			t.st.Classes[cls.Name].Methods[name] = f
		}
	}
}
