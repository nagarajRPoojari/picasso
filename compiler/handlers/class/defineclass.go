package class

import (
	"github.com/nagarajRPoojari/x-lang/ast"
	funcs "github.com/nagarajRPoojari/x-lang/compiler/handlers/func"
)

// defineClass similar to declareClass but does function concrete declaration
func (t *ClassHandler) DefineClass(cls ast.ClassDeclarationStatement) {
	for _, stI := range cls.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			funcs.FuncHandlerInst.DefineFunc(cls.Name, &st)
		}
	}
}
