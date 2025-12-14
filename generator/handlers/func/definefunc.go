package funcs

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/block"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/constants"
	tf "github.com/nagarajRPoojari/x-lang/generator/type"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
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

	bh := bc.NewBlockHolder(bc.VarBlock{Block: f.NewBlock("")}, f.NewBlock(""))

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
			cls := &tf.Class{
				Name: className,
				UDT:  clsMeta.UDT.(*types.PointerType),
			}
			cls.Update(bh, p)
			t.st.Vars.AddNewVar(p.LocalName, cls)
		}
	}

	old := bh.N
	block.BlockHandlerInst.ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)
	if fn.ReturnType == nil {
		bh.N.NewRet(nil)
	}
}

func (t *FuncHandler) DefineMainFunc(fn *ast.FunctionDefinitionStatement, avoid map[string]struct{}) {
	// new level for function block
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	var f *ir.Func = t.st.MainFunc
	bh := bc.NewBlockHolder(bc.VarBlock{Block: f.NewBlock(constants.ENTRY)}, f.NewBlock(""))
	t.Init(bh)
	bh.N.NewCall(t.st.CI.Funcs[c.FUNC_RUNTIME_INIT])

	if len(fn.Parameters) != 0 {
		errorutils.Abort(errorutils.MainFuncError, "parameters are not allowed in main function")
	}

	old := bh.N
	block.BlockHandlerInst.ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)

	nullPtr := constant.NewNull((types.NewPointer(types.I8)))

	// Return it
	bh.N.NewRet(nullPtr)
}

func (t *FuncHandler) Init(block *bc.BlockHolder) {
	tps := []string{"int64", "int32", "int16", "int8", "string"}
	for _, tp := range tps {
		t.initTypes(block, tp)
	}
}

func (t *FuncHandler) initTypes(block *bc.BlockHolder, s string) {
	t.st.Vars.RegisterTypeHolders(block, s, t.st.TypeHandler.BuildVar(block, tf.NewType(s), nil))
}
