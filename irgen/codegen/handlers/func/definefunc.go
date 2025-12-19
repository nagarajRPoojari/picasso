package funcs

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/block"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/constants"
	tf "github.com/nagarajRPoojari/niyama/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

// DefineFunc generates the concrete LLVM IR body for a class method or constructor.
// It initializes the function's entry blocks, populates the local symbol table with
// parameters (including the implicit 'this' pointer), and delegates statement
// lowering to the BlockHandler.
//
// Technical Logic:
//   - Scope Management: Opens a fresh 'Func' level variable scope, ensuring local
//     variables do not leak between function boundaries.
//   - Entry Point Logic: Creates two initial blocks—one for stack allocations (Alloca)
//     and one for the actual execution logic (the entry block).
//   - Constructor Initialization: If the method is identified as a constructor
//     (name matches class name), it triggers 'initTypes' to set up class-specific defaults.
//   - Parameter Binding: Iterates through LLVM formal parameters to register them
//     in the symbol table. It distinguishes between standard user-defined parameters
//     and the implicit 'this' pointer, which is wrapped in a tf.Class container.
//   - Terminal Handling: Automatically injects a 'void' return if no explicit
func (t *FuncHandler) DefineFunc(className string, fn *ast.FunctionDefinitionStatement, avoid map[string]struct{}) {
	// new level for function block
	t.st.Vars.AddFunc()
	defer t.st.Vars.RemoveFunc()

	name := fmt.Sprintf("%s.%s", className, fn.Name)
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

// DefineMainFunc generates the entry point for the executable. Unlike standard
// class methods, the main function is responsible for bootstrapping the
// Niyama runtime and initializing global state before executing the user's code.
//
// Technical Logic:
//   - Runtime Bootstrapping: Injects a call to the internal runtime initialization
//     function (e.g., GC setup, thread pool init) as the first action.
//   - Entry Block Management: Creates a dedicated 'entry' block for variable
//     allocations to ensure all stack pointers are resolved at the start of the function.
//   - Signature Enforcement: Strictly validates that the main function
//     contains no parameters, ensuring compliance with the language spec.
//   - Exit Strategy: Automatically returns a null pointer (i8*) upon completion,
//     serving as the standard exit signal for the host environment.
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
	// initTypes are to create special variables of named primitive types.
	tps := []string{"int64", "int32", "int16", "int8", "string"}
	for _, tp := range tps {
		t.initTypes(block, tp)
	}
}

func (t *FuncHandler) initTypes(block *bc.BlockHolder, s string) {
	t.st.Vars.RegisterTypeHolders(block, s, t.st.TypeHandler.BuildVar(block, tf.NewType(s), nil))
}
