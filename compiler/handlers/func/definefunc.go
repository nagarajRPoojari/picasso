package funcs

import (
	"github.com/llir/llvm/ir"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/block"
	"github.com/nagarajRPoojari/x-lang/compiler/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
)

// defineFunc does concrete function declaration
func (t *FuncHandler) DefineFunc(className string, fn *ast.FunctionDefinitionStatement, avoid map[string]struct{}) {
	// new level for function block
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	name := t.st.IdentifierBuilder.Attach(className, fn.Name)
	var f *ir.Func
	f = t.st.Classes[className].Methods[name]
	if _, ok := avoid[fn.Name]; ok {
		return
	}

	bh := tf.BlockHolder{V: tf.VarBlock{Block: f.NewBlock("")}, N: f.NewBlock(constants.ENTRY)}

	if className == fn.Name {
		t.initTypes(bh, className)
	}

	if name == constants.MAIN && len(fn.Parameters) != 0 {
		errorutils.Abort(errorutils.MainFuncError, "parameters are not allowed in main function")
	}

	for i, p := range f.Params {
		if i < len(fn.Parameters) {
			pt := fn.Parameters[i].Type
			paramType := tf.NewType(pt.Get(), pt.GetUnderlyingType())
			t.st.Vars.AddNewVar(p.LocalName, t.st.TypeHandler.BuildVar(bh, paramType, p))
		} else {
			clsMeta := t.st.Classes[className]
			if clsMeta == nil {
				errorutils.Abort(errorutils.UnknownClass, className)
			}
			t.st.Vars.AddNewVar(p.LocalName, &tf.Class{
				Name: className,
				UDT:  clsMeta.UDT,
				Ptr:  p,
			})
		}
	}

	old := bh.N
	bh = block.BlockHandlerInst.ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)
	if fn.ReturnType == nil {
		bh.N.NewRet(nil)
	}
}

func (t *FuncHandler) DefineMainFunc(fn *ast.FunctionDefinitionStatement, avoid map[string]struct{}) {
	// new level for function block
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	var f *ir.Func
	f = t.st.MainFunc
	bh := tf.BlockHolder{V: tf.VarBlock{Block: f.NewBlock(constants.ENTRY)}, N: f.NewBlock("")}
	t.Init(bh)
	bh.N.NewCall(t.st.GC.Init())

	if len(fn.Parameters) != 0 {
		errorutils.Abort(errorutils.MainFuncError, "parameters are not allowed in main function")
	}

	old := bh.N
	bh = block.BlockHandlerInst.ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)
	if fn.ReturnType == nil {
		bh.N.NewRet(nil)
	}
}

func (t *FuncHandler) Init(block tf.BlockHolder) {
	tps := []string{"int64", "int32", "int16", "int8", "string"}
	for _, tp := range tps {
		t.initTypes(block, tp)
	}
}

func (t *FuncHandler) initTypes(block tf.BlockHolder, s string) {
	t.st.Vars.RegisterTypeHolders(block, s, t.st.TypeHandler.BuildVar(block, tf.NewType(s), nil))
}
