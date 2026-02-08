package funcs

import (
	"fmt"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/block"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/expression"
	tf "github.com/nagarajRPoojari/picasso/irgen/codegen/type"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
)

func (t *FuncHandler) DefineConstructor(className string, instance *tf.Class, bh *bc.BlockHolder, fn *ast.FunctionDefinitionStatement, avoid map[string]struct{}) *tf.Class {
	aliasClsName := className

	if _, ok := t.st.Interfaces[aliasClsName]; ok {
		errorutils.Abort(errorutils.InterfaceInstantiationError, aliasClsName)
	}

	classMeta := t.st.Classes[aliasClsName]
	if classMeta == nil {
		errorutils.Abort(errorutils.UnknownClass, aliasClsName)
	}

	clsNameSplit := strings.Split(aliasClsName, ".")
	moduleName := strings.Join(clsNameSplit[:len(clsNameSplit)-1], ".")

	if classMeta.Internal && moduleName != t.st.ModuleName {
		errorutils.Abort(errorutils.ClassNotAccessible, aliasClsName)
	}

	// tf.NewClass allocates memory for class instance in heap internally.
	// & holds heap pointer in a stack slot.
	structType := classMeta.StructType()
	meta := t.st.Classes[aliasClsName]

	for name, index := range meta.FieldIndexMap {
		// if field is not found in meta.VarAST indicating func type, update instance
		// to point to the function. future function calls on that instance will directly
		// refer to this pointed function.
		exp, ok := meta.VarAST[name]
		if !ok {
			f := t.st.Classes[aliasClsName].Methods[name]
			fieldType := t.st.Classes[aliasClsName].StructType().Fields[index]
			instance.UpdateField(bh, t.st.TypeHandler, index, f, fieldType)
			continue
		}

		// handling initialized & uninitialized variables.
		fieldType := structType.Fields[index]
		var v tf.Var
		if exp.AssignedValue == nil {
			// @todo: this need to be verified
			var init value.Value

			// atomic data types are special class types & are not expected to be initialized with
			// new keyword. e.g, say x: atomic int; should do the instantiaion job though it is just
			// a declaration. therefore instantiate with NewClass.
			if exp.ExplicitType.IsAtomic() {
				meta := t.st.Classes[exp.ExplicitType.Get()]
				v = tf.NewClass(bh, exp.ExplicitType.Get(), meta.UDT)
			} else {

				// remaining vars without assignedvalues holds its corresponding zero values.
				// @todo: list zero values for all data types somewehere in docs to look at.
				v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get(), exp.ExplicitType.GetUnderlyingType()), init)
			}
		} else {
			v = t.m.GetExpressionHandler().(*expression.ExpressionHandler).ProcessExpression(bh, exp.AssignedValue)

			// data types other than array, like primitives, object types are typecasted implicitly
			// before assignment.
			if v.NativeTypeString() != constants.ARRAY {
				casted := t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
				v = t.st.TypeHandler.BuildVar(bh, tf.NewType(exp.ExplicitType.Get()), casted)
			} else {
				// no need to cast array type, but do a base type check.
				t.st.TypeHandler.ImplicitTypeCast(bh, exp.ExplicitType.Get(), v.Load(bh))
			}
		}
		instance.UpdateField(bh, t.st.TypeHandler, index, v.Load(bh), fieldType)
		// t.st.Vars.AddNewVar(exp.Identifier, v)
	}

	return instance
}

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
		// @todo: not sure about its placement here, could have been handled
		// while declaring
		errorutils.Abort(errorutils.MethodRedeclaration, fn.Name)
		return
	}

	clsAliasNameSplit := strings.Split(className, ".")
	clsAliasName := clsAliasNameSplit[len(clsAliasNameSplit)-1]
	var instance value.Value
	if clsAliasName == fn.Name {
		// constructor: init Types
		// @todo: basic checks about constructor
		// t.initTypes(bh, className)
	}
	bh := bc.NewBlockHolder(bc.VarBlock{Block: f.NewBlock("")}, f.NewBlock(""))
	old := bh.N

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
			if clsAliasName == fn.Name {
				t.DefineConstructor(className, cls, bh, fn, avoid)
				instance = cls.Load(bh)
			}
			t.st.Vars.AddNewVar(p.LocalName, cls)
		}
	}

	t.m.GetBlockHandler().(*block.BlockHandler).ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)
	if clsAliasName == fn.Name {
		bh.N.NewRet(instance)
		return
	}

	if fn.ReturnType == nil {
		bh.N.NewRet(nil)
	}
}

// DefineMainFunc generates the entry point for the executable. Unlike standard
// class methods, the main function is responsible for bootstrapping the
// Picasso runtime and initializing global state before executing the user's code.
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
	// t.Init(bh)
	bh.N.NewCall(t.st.CI.Funcs[c.FUNC_RUNTIME_INIT])

	if len(fn.Parameters) != 0 {
		errorutils.Abort(errorutils.InvalidMainMethodSignature, "parameters are not allowed in main function")
	}
	if fn.ReturnType != nil {
		errorutils.Abort(errorutils.InvalidMainMethodSignature, "expected no return type for main functions")
	}

	old := bh.N
	t.m.GetBlockHandler().(*block.BlockHandler).ProcessBlock(f, bh, fn.Body)
	bh.V.NewBr(old)

	nullPtr := constant.NewNull((types.NewPointer(types.I8)))

	// Return it
	bh.N.NewRet(nullPtr)
}
