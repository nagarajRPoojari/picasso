package compiler

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
	function "github.com/nagarajRPoojari/x-lang/compiler/libs/func"
	tf "github.com/nagarajRPoojari/x-lang/compiler/type"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// LLVM parses abstract syntax tree to generate llvm IR
type LLVM struct {

	// llvm module
	module *ir.Module

	typeHandler *tf.TypeHandler
	// utility tool to build consistent identifier names
	identifierBuilder *IdentifierBuilder

	// all global vars
	vars map[string]tf.Var
	// all methods including class methods & top level functions
	methods map[string]*ir.Func
	// custom classes defined by user
	classes map[string]*tf.MetaClass

	// global string counter
	// @todo: move this to separate string module
	strCounter int

	LibMethods map[string]function.Func
}

func NewLLVM() *LLVM {
	m := ir.NewModule()
	i := &LLVM{
		module:            m,
		typeHandler:       tf.NewTypeHandler(),
		vars:              make(map[string]tf.Var),
		methods:           make(map[string]*ir.Func),
		classes:           make(map[string]*tf.MetaClass),
		identifierBuilder: NewIdentifierBuilder(MAIN),
		LibMethods:        make(map[string]function.Func),
	}
	return i
}

func (t *LLVM) Dump(file string) {
	f, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(t.module.String())
}

func (t *LLVM) ParseAST(tree *ast.BlockStatement) {
	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ImportStatement:
			t.importer(t.LibMethods, st.Name)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			t.predeclareClass(st)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			t.defineClassVars(st)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			t.declareFunctions(st)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			t.defineClass(st)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			if st.Name == MAIN {
				f := t.module.NewFunc(MAIN, types.I32)
				t.methods[MAIN] = f
				t.defineFunc("", &st)
			}
		}
	}
}

// predeclareClass creates an opaque struct for all classes defined by user
// and registers it with typehandler for identfying forward declaration
func (t *LLVM) predeclareClass(class ast.ClassDeclarationStatement) {
	if _, ok := t.classes[class.Name]; ok {
		return
	}
	udt := types.NewStruct() // opaque
	t.module.NewTypeDef(class.Name, udt)
	mc := &tf.MetaClass{
		VarIndexMap: make(map[string]int),
		VarAST:      make(map[string]*ast.VariableDeclarationStatement),
		Methods:     make(map[string]*ir.Func),
		UDT:         udt,
	}
	t.classes[class.Name] = mc
	t.typeHandler.Register(mc)
}

// defineClassVars stores corresponding ast for all var declaration
// which will be used to instantiate them on constructor call, i.e, new MyClass()
func (t *LLVM) defineClassVars(class ast.ClassDeclarationStatement) {
	mc := t.classes[class.Name]
	fieldTypes := make([]types.Type, 0)

	// map each fields with corresponding udt struct index
	i := 0
	for _, stI := range class.Body {
		if st, ok := stI.(ast.VariableDeclarationStatement); ok {
			fqName := t.identifierBuilder.Attach(class.Name, st.Identifier)
			mc.VarIndexMap[fqName] = i
			mc.VarAST[fqName] = &st

			fieldType := t.typeHandler.GetLLVMType(tf.Type(st.ExplicitType.Get()))
			fieldTypes = append(fieldTypes, fieldType)
			i++
		}
	}

	// update opaque udt with concreate struct
	mc.UDT.(*types.StructType).Fields = fieldTypes
}

// declareFunctions loops over all functions inside Class & creates
// a header declaration
func (t *LLVM) declareFunctions(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.typeHandler.GetLLVMType(tf.Type(p.Type.Get()))))
			}

			// at the end pass `this` parameter representing current object
			udt := t.classes[class.Name].UDT
			thisParamType := types.NewPointer(udt)
			params = append(params, ir.NewParam(THIS, thisParamType))

			name := t.identifierBuilder.Attach(class.Name, st.Name)

			// @todo: handle no return type case
			var retType types.Type
			if st.ReturnType != nil {
				retType = t.typeHandler.GetLLVMType(tf.Type(st.ReturnType.Get()))
			} else {
				retType = t.typeHandler.GetLLVMType(tf.Type(tf.NULL))
			}
			f := t.module.NewFunc(name, retType, params...)
			t.methods[name] = f
			t.classes[class.Name].Methods[name] = f
		}
	}
}

// defineClass similar to declareClass but does function concrete declaration
func (t *LLVM) defineClass(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			t.defineFunc(class.Name, &st)
		}
	}
}

// defineFunc does concrete function declaration
func (t *LLVM) defineFunc(className string, fn *ast.FunctionDeclarationStatement) {
	vars := make(map[string]tf.Var, 0)
	name := t.identifierBuilder.Attach(className, fn.Name)
	if className == "" { // indicates classless function: main
		name = fn.Name
	}
	f := t.methods[name]
	entry := f.NewBlock(ENTRY)

	if name == MAIN && len(fn.Parameters) != 0 {
		errorsx.PanicCompilationError("parameters are not allowed in main function")
	}

	for i, p := range f.Params {
		if i < len(fn.Parameters) {
			paramType := tf.Type(fn.Parameters[i].Type.Get())
			vars[p.LocalName] = t.typeHandler.BuildVar(entry, paramType, p)
			continue
		}

		clsMeta := t.classes[className]
		if clsMeta == nil {
			errorsx.PanicCompilationError("defineFunc: unknown class when binding this")
		}
		vars[p.LocalName] = &tf.Class{
			Name: className,
			UDT:  clsMeta.UDT,
			Ptr:  p,
		}
		break
	}

	for _, stI := range fn.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			var v tf.Var
			if st.AssignedValue == nil {
				v = t.typeHandler.BuildDefaultVar(entry, tf.Type(st.ExplicitType.Get()))
				vars[st.Identifier] = v
			} else {
				v = t.processExpression(entry, vars, st.AssignedValue)
				if _, ok := vars[st.Identifier]; ok {
					errorsx.PanicCompilationError("variable already exists")
				}
				casted := t.typeHandler.CastToType(entry, st.ExplicitType.Get(), v.Load(entry))
				v = t.typeHandler.BuildVar(entry, tf.Type(st.ExplicitType.Get()), casted)
				vars[st.Identifier] = v
			}

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:

				assigneeExp, _ := exp.Assignee.(ast.SymbolExpression)
				assignee := assigneeExp.Value
				v, ok := vars[assignee]
				if !ok {
					globalV, ok := t.vars[assignee]
					if !ok {
						panic(fmt.Sprintf("undefined: %s", assignee))
					} else {
						v = globalV
					}
				}
				rhs := t.processExpression(entry, vars, exp.AssignedValue)
				typeName := v.Type().Name()
				if typeName == "" {
					typeName = v.Type().String()
				}
				casted := t.typeHandler.CastToType(entry, typeName, rhs.Load(entry))
				c := t.typeHandler.BuildVar(entry, tf.Type(typeName), casted)
				v.Update(entry, c.Load(entry))

			case ast.CallExpression:
				// will be routed CallExpression
				t.processExpression(entry, vars, exp)
			default:
				errorsx.PanicCompilationError("invalid statement")
			}
		case ast.ReturnStatement:
			v := t.processExpression(entry, vars, st.Value.Expression)
			entry.NewRet(t.typeHandler.CastToType(entry, fn.ReturnType.Get(), v.Load(entry)))
		}
	}

	// @todo: remove this debug statement
	// if name == MAIN {
	// 	x := vars["z"].Load(entry)
	// 	t.print(entry, "return z = %s", x)
	// }

	if fn.ReturnType == nil {
		entry.NewRet(constant.NewNull(types.NewPointer(types.NewStruct())))
	}
}

// processExpression handles binary expressions, function calls, member operations etc..
func (t *LLVM) processExpression(block *ir.Block, vars map[string]tf.Var, expI ast.Expression) tf.Var {
	if expI == nil {
		return tf.NewNullVar(types.NewPointer(types.NewStruct()))
	}

	switch ex := expI.(type) {

	case ast.SymbolExpression:
		// search for variable locally then gloablly
		if v, ok := vars[ex.Value]; ok {
			return v
		}
		if gv, ok := t.vars[ex.Value]; ok {
			return gv
		}
		errorsx.PanicCompilationError(fmt.Sprintf("undefined var: %s", ex.Value))

	case ast.NumberExpression:
		// produce a runtime mutable var for the literal (double)
		// by default number will be wrapped up with float64
		return t.typeHandler.BuildVar(block, tf.FLOAT64, constant.NewFloat(types.Double, ex.Value))

	case ast.StringExpression:
		formatStr := ex.Value
		strConst := constant.NewCharArrayFromString(formatStr + "\x00")
		global := t.module.NewGlobalDef("", strConst)

		gep := block.NewGetElementPtr(
			global.ContentType,
			global,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 0),
		)

		return tf.NewString(block, gep)

	case ast.NewExpression:
		// @todo: make constructor call
		meth := ex.Instantiation.Method.(ast.SymbolExpression)
		classMeta := t.classes[meth.Value]
		if classMeta == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("unknown class: %s", meth.Value))
		}

		instance := tf.NewClass(block, meth.Value, classMeta.UDT)
		structType := classMeta.UDT.(*types.StructType)
		meta := t.classes[meth.Value]

		for name, index := range meta.VarIndexMap {
			exp := meta.VarAST[name]
			x := t.processExpression(block, vars, exp.AssignedValue)

			fieldType := structType.Fields[index]
			instance.UpdateField(block, index, x.Load(block), fieldType)
		}
		return instance

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar := t.processExpression(block, vars, ex.Member)
		if baseVar == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("nil base in member expression: %v %v", ex.Member, vars))
		}

		// Base must be a class instance
		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("member access base is not a class instance, got %T", baseVar))
		}

		// Get metadata for base class
		classMeta, ok := t.classes[cls.Name]
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("unknown class metadata: %s", cls.Name))
		}

		// Compute field name in identifier map
		fieldID := t.identifierBuilder.Attach(cls.Name, ex.Property)
		idx, ok := classMeta.VarIndexMap[fieldID]
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("unknown field %s on class %s", ex.Property, cls.Name))
		}

		// Get field type from struct UDT
		st, ok := classMeta.UDT.(*types.StructType)
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("class %s does not have a struct UDT", cls.Name))
		}
		fieldType := st.Fields[idx]

		// Get pointer to the field
		fieldPtr := cls.FieldPtr(block, idx)

		// Determine the class name if the field is a struct
		getClassName := func(tt types.Type) string {
			for cname, meta := range t.classes {
				if meta.UDT == tt {
					return cname
				}
			}
			for cname, meta := range t.typeHandler.Udts {
				if meta.UDT == tt {
					return cname
				}
			}
			return ""
		}

		// Wrap into appropriate Var
		switch ft := fieldType.(type) {
		case *types.IntType:
			switch ft.BitSize {
			case 1:
				return &tf.Boolean{NativeType: types.I1, Value: fieldPtr}
			case 8:
				return &tf.Int8{NativeType: types.I8, Value: fieldPtr}
			case 16:
				return &tf.Int16{NativeType: types.I16, Value: fieldPtr}
			case 32:
				return &tf.Int32{NativeType: types.I32, Value: fieldPtr}
			case 64:
				return &tf.Int64{NativeType: types.I64, Value: fieldPtr}
			default:
				panic(fmt.Sprintf("unsupported int size %d", ft.BitSize))
			}

		case *types.FloatType:
			switch ft.Kind {
			case types.FloatKindHalf:
				return &tf.Float16{NativeType: types.Half, Value: fieldPtr}
			case types.FloatKindFloat:
				return &tf.Float32{NativeType: types.Float, Value: fieldPtr}
			case types.FloatKindDouble:
				return &tf.Float64{NativeType: types.Double, Value: fieldPtr}
			default:
				panic(fmt.Sprintf("unsupported float kind %v", ft.Kind))
			}

		case *types.StructType:
			return &tf.Class{
				Name: getClassName(fieldType),
				UDT:  fieldType,
				Ptr:  fieldPtr,
			}

		default:
			errorsx.PanicCompilationError(fmt.Sprintf("unsupported field type %T in member expression", fieldType))
		}

	case ast.CallExpression:
		return t.handleCallExpression(block, vars, ex)

	case ast.BinaryExpression:
		left := t.processExpression(block, vars, ex.Left)
		right := t.processExpression(block, vars, ex.Right)
		if left == nil || right == nil {
			errorsx.PanicCompilationError("nil operand in binary expression")
		}
		lv := left.Load(block)
		rv := right.Load(block)

		f := &tf.Float64{}
		lvf, err := f.Cast(block, lv)
		if err != nil {
			errorsx.PanicCompilationError((fmt.Sprintf("failed to cast %s to float", lv)))
		}
		rvf, err := f.Cast(block, rv)
		if err != nil {
			errorsx.PanicCompilationError((fmt.Sprintf("failed to cast %s to float", rv)))
		}

		var res value.Value
		switch ex.Operator.Value {
		case "+":
			res = block.NewFAdd(lvf, rvf)
		case "-":
			res = block.NewFSub(lvf, rvf)
		case "*":
			res = block.NewFMul(lvf, rvf)
		case "/":
			res = block.NewFDiv(lvf, rvf)
		default:
			panic(fmt.Sprintf("unsupported operator: %s", ex.Operator.Value))
		}

		ptr := block.NewAlloca(types.Double)
		block.NewStore(res, ptr)
		return &tf.Float64{NativeType: types.Double, Value: ptr}
	}

	return nil
}

func (t *LLVM) handleCallExpression(block *ir.Block, vars map[string]tf.Var, ex ast.CallExpression) tf.Var {
	// check if imported modules
	if m, ok := ex.Method.(ast.MemberExpression); ok {
		fName := fmt.Sprintf("%s.%s", m.Member.(ast.SymbolExpression).Value, m.Property)
		if f, ok := t.LibMethods[fName]; ok {
			args := make([]tf.Var, 0)
			for _, v := range ex.Arguments {
				args = append(args, t.processExpression(block, vars, v))
			}
			return f(t.typeHandler, t.module, block, args)
		}
	}

	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		errorsx.PanicCompilationError("method call should be on instance")

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar := t.processExpression(block, vars, m.Member)
		if baseVar == nil {
			errorsx.PanicCompilationError("handleCallExpression: nil baseVar for member expression")
		}

		cls, ok := baseVar.(*tf.Class)
		if !ok {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: member access base is not Class (got %T)", baseVar))
		}
		if cls == nil || cls.Ptr == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: class or class.Ptr is nil for class %v", cls))
		}

		classMeta := t.classes[cls.Name]
		if classMeta == nil {
			errorsx.PanicCompilationError("handleCallExpression: unknown class metadata: " + cls.Name)
		}

		methodKey := t.identifierBuilder.Attach(cls.Name, m.Property)
		fn, ok := classMeta.Methods[methodKey]
		if !ok || fn == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: unknown method %s on class %s", m.Property, cls.Name))
		}

		// Build args for the user parameters (do not append `this` yet)
		args := make([]value.Value, 0, len(ex.Arguments)+1)
		for i, argExp := range ex.Arguments {
			v := t.processExpression(block, vars, argExp)
			if v == nil {
				errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: nil arg %d for %s.%s", i, cls.Name, m.Property))
			}
			raw := v.Load(block)
			if raw == nil {
				errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: loaded nil arg %d for %s.%s", i, cls.Name, m.Property))
			}

			// If the callee expects a certain param type, cast to it
			expected := fn.Sig.Params[i]
			raw = t.typeHandler.CastToType(block, expected.Name(), raw)
			if raw == nil {
				errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: CastToType returned nil for arg %d -> %s", i, expected.String()))
			}
			args = append(args, raw)
		}

		// Pass `this` as a pointer-to-struct (Slot returns pointer)
		thisPtr := cls.Slot()
		if thisPtr == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: this pointer is nil for %s", cls.Name))
		}

		// Check function expected param count: we declared 'this' last when creating fn,
		// adjust order according to how the function was declared.
		args = append(args, thisPtr)

		// ensure function is non-nil
		if fn == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: function pointer nil for %s.%s", cls.Name, m.Property))
		}

		// Now call
		ret := block.NewCall(fn, args...)
		if fn.Sig.RetType == types.Void {
			return nil
		}

		// allocate slot for return, store and wrap
		if ret == nil {
			errorsx.PanicCompilationError(fmt.Sprintf("handleCallExpression: call returned nil for %s.%s", cls.Name, m.Property))
		}
		slot := block.NewAlloca(fn.Sig.RetType)
		if slot == nil {
			errorsx.PanicCompilationError("handleCallExpression: failed to alloca for return")
		}
		block.NewStore(ret, slot)

		return t.wrapReturn(slot, fn.Sig.RetType, cls.Name)
	}
	return nil
}

func (t *LLVM) wrapReturn(slot *ir.InstAlloca, rt types.Type, className string) tf.Var {
	switch v := rt.(type) {
	case *types.IntType:
		switch v.BitSize {
		case 1:
			return &tf.Boolean{NativeType: types.I1, Value: slot}
		case 8:
			return &tf.Int8{NativeType: types.I8, Value: slot}
		case 16:
			return &tf.Int16{NativeType: types.I16, Value: slot}
		case 32:
			return &tf.Int32{NativeType: types.I32, Value: slot}
		case 64:
			return &tf.Int64{NativeType: types.I64, Value: slot}
		}
	case *types.FloatType:
		switch v.Kind {
		case types.FloatKindFloat:
			return &tf.Float32{NativeType: types.Float, Value: slot}
		case types.FloatKindDouble:
			return &tf.Float64{NativeType: types.Double, Value: slot}
		}
	case *types.StructType:
		return &tf.Class{Name: v.Name(), UDT: v, Ptr: slot}
	}
	return &tf.Class{Name: className, UDT: rt, Ptr: slot}
}
