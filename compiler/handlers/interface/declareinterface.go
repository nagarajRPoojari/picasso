package _interface

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// predeclareClass creates an opaque struct for all classes defined by user
// and registers it with typehandler for identfying forward declaration
func (t *InterfaceHandler) PredeclareInterface(cls ast.ClassDeclarationStatement) {
	if _, ok := t.st.Classes[cls.Name]; ok {
		return
	}
	udt := types.NewStruct() // opaque
	t.st.Module.NewTypeDef(cls.Name, udt)
	mc := &tf.MetaClass{
		VarIndexMap: make(map[string]int),
		VarAST:      make(map[string]*ast.VariableDeclarationStatement),
		Methods:     make(map[string]*ir.Func),
		UDT:         types.NewPointer(udt),
	}
	t.st.Classes[cls.Name] = mc
	t.st.TypeHandler.Register(cls.Name, mc)
}
