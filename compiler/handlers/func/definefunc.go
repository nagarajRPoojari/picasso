package funcs

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/block"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// defineFunc does concrete function declaration
func (t *FuncHandler) DefineFunc(className string, fn *ast.FunctionDefinitionStatement) {
	// new level for function block
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	name := t.st.IdentifierBuilder.Attach(className, fn.Name)
	var f *ir.Func
	if className == "" { // indicates classless function: main
		name = fn.Name
		f = t.st.MainFunc
	} else {
		f = t.st.Classes[className].Methods[name]
	}
	entry := f.NewBlock(constants.ENTRY)

	if className == "" {
		entry.NewCall(t.st.GC.Init())
	}

	if name == constants.MAIN && len(fn.Parameters) != 0 {
		errorsx.PanicCompilationError("parameters are not allowed in main function")
	}

	for i, p := range f.Params {
		if i < len(fn.Parameters) {
			paramType := tf.Type(fn.Parameters[i].Type.Get())
			t.st.Vars.AddNewVar(p.LocalName, t.st.TypeHandler.BuildVar(entry, paramType, p))
			continue
		}

		clsMeta := t.st.Classes[className]
		if clsMeta == nil {
			errorsx.PanicCompilationError("defineFunc: unknown class when binding this")
		}
		t.st.Vars.AddNewVar(p.LocalName, &tf.Class{
			Name: className,
			UDT:  clsMeta.UDT,
			Ptr:  p,
		})
		break
	}

	entry = block.BlockHandlerInst.ProcessBlock(f, entry, fn.Body)

	if fn.ReturnType == nil {
		entry.NewRet(nil)
	}
}
