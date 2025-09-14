package compiler

import (
	"fmt"
	"os"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/ast"
)

type MetaClass struct {
	varIndexMap map[string]int
	varAST      map[string]*ast.VariableDeclarationStatement

	methods map[string]*ir.Func

	udt types.Type
}

func NewMetaClass() *MetaClass {
	return &MetaClass{
		varIndexMap: make(map[string]int),
		varAST:      make(map[string]*ast.VariableDeclarationStatement),
		methods:     make(map[string]*ir.Func),
	}
}
func (m *MetaClass) FieldType(idx int) types.Type {
	return m.udt.(*types.StructType).Fields[idx]
}

type LLVM struct {
	module *ir.Module

	typeHandler       *TypeHandler
	identifierBuilder *IdentifierBuilder

	vars    map[string]Var
	methods map[string]*ir.Func
	classes map[string]*MetaClass

	classLookUp map[string]struct{}

	strCounter int
}

func NewLLVM() *LLVM {
	m := ir.NewModule()
	i := &LLVM{
		module:            m,
		typeHandler:       NewTypeHandler(),
		vars:              make(map[string]Var),
		methods:           make(map[string]*ir.Func),
		classes:           make(map[string]*MetaClass),
		identifierBuilder: NewIdentifierBuilder(MAIN),
	}
	return i
}

func (t *LLVM) Dump() {
	f, err := os.Create("bin/output.ll")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(t.module.String())
}

func (t *LLVM) ParseAST(tree *ast.BlockStatement) {
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

func (t *LLVM) predeclareClass(class ast.ClassDeclarationStatement) {
	if _, ok := t.classes[class.Name]; ok {
		return
	}
	udt := types.NewStruct() // opaque
	t.module.NewTypeDef(class.Name, udt)
	mc := &MetaClass{
		varIndexMap: make(map[string]int),
		varAST:      make(map[string]*ast.VariableDeclarationStatement),
		methods:     make(map[string]*ir.Func),
		udt:         udt,
	}
	t.classes[class.Name] = mc
	t.typeHandler.Register(mc)
}

func (t *LLVM) defineClassVars(class ast.ClassDeclarationStatement) {
	mc := t.classes[class.Name]
	fieldTypes := make([]types.Type, 0)

	i := 0
	for _, stI := range class.Body {
		if st, ok := stI.(ast.VariableDeclarationStatement); ok {
			fqName := t.identifierBuilder.Attach(class.Name, st.Identifier)
			mc.varIndexMap[fqName] = i
			mc.varAST[fqName] = &st

			fieldType := t.typeHandler.GetLLVMType(Type(st.ExplicitType.Get()))
			fieldTypes = append(fieldTypes, fieldType)
			i++
		}
	}

	mc.udt.(*types.StructType).Fields = fieldTypes
}

func (t *LLVM) declareFunctions(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.typeHandler.GetLLVMType(Type(p.Type.Get()))))
			}

			udt := t.classes[class.Name].udt
			thisParamType := types.NewPointer(udt)
			params = append(params, ir.NewParam(THIS, thisParamType))

			name := t.identifierBuilder.Attach(class.Name, st.Name)
			retType := t.typeHandler.GetLLVMType(Type(st.ReturnType.Get()))
			f := t.module.NewFunc(name, retType, params...)
			t.methods[name] = f
			t.classes[class.Name].methods[name] = f
		}
	}
}

func (t *LLVM) defineClass(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			t.defineFunc(class.Name, &st)
		}
	}
}

func (t *LLVM) defineFunc(className string, fn *ast.FunctionDeclarationStatement) {
	vars := make(map[string]Var, 0)
	name := t.identifierBuilder.Attach(className, fn.Name)
	if className == "" { // indicates classless function: main
		name = fn.Name
	}
	f := t.methods[name]
	entry := f.NewBlock(ENTRY)

	if name == MAIN && len(fn.Parameters) != 0 {
		panic("parameters are not allowed in main function")
	}

	for i, p := range f.Params {
		if i < len(fn.Parameters) {
			paramType := Type(fn.Parameters[i].Type.Get())
			vars[p.LocalName] = t.typeHandler.BuildVar(entry, paramType, p)
			continue
		}

		clsMeta := t.classes[className]
		if clsMeta == nil {
			panic("defineFunc: unknown class when binding this")
		}
		vars[p.LocalName] = &Class{
			Name: className,
			UDT:  clsMeta.udt,
			Ptr:  p,
		}
		break
	}

	for _, stI := range fn.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			v := t.processExpression(entry, vars, st.AssignedValue)
			if _, ok := vars[st.Identifier]; ok {
				panic("variable already exists")
			}

			if st.ExplicitType != nil {
				casted := t.typeHandler.CastToType(entry, st.ExplicitType.Get(), v.Load(entry))
				v = t.typeHandler.BuildVar(entry, Type(st.ExplicitType.Get()), casted)
			}

			vars[st.Identifier] = v

		case ast.ExpressionStatement:
			switch exp := st.Expression.(type) {
			case ast.AssignmentExpression:
				// @todo: type casting Var
				// @todo: handle member assignments. e.g, this.x = 100;
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
				fmt.Println("assigning ", rhs)
				v.Update(entry, rhs.Load(entry))
			case ast.CallExpression:
				// will be routed CallExpression
				t.processExpression(entry, vars, exp)
			default:
				panic("invalid statement")
			}
		case ast.ReturnStatement:
			v := t.processExpression(entry, vars, st.Value.Expression)
			entry.NewRet(t.typeHandler.CastToType(entry, fn.ReturnType.Get(), v.Load(entry)))
		}
	}

	if name == MAIN {
		x := vars["z"].Load(entry)
		t.print(entry, "return z = %f %d", x, x)
	}

}

// processExpression: handles symbols, numbers, new, member, call, binary etc.
// Supports chaining like a.b().c().d()
func (t *LLVM) processExpression(block *ir.Block, vars map[string]Var, expI ast.Expression) Var {
	switch ex := expI.(type) {

	case ast.SymbolExpression:
		if v, ok := vars[ex.Value]; ok {
			return v
		}
		if gv, ok := t.vars[ex.Value]; ok {
			return gv
		}
		panic(fmt.Sprintf("undefined: %s", ex.Value))

	case ast.NumberExpression:
		// produce a runtime mutable var for the literal (double)
		return t.typeHandler.GetPrimitiveVar(block, FLOAT64, constant.NewFloat(types.Double, ex.Value))

	case ast.NewExpression:
		meth := ex.Instantiation.Method.(ast.SymbolExpression)
		classMeta := t.classes[meth.Value]
		if classMeta == nil {
			panic(fmt.Sprintf("unknown class: %s", meth.Value))
		}

		instance := NewClass(block, meth.Value, classMeta.udt)
		structType := classMeta.udt.(*types.StructType)
		meta := t.classes[meth.Value]

		for name, index := range meta.varIndexMap {
			exp := meta.varAST[name]
			x := t.processExpression(block, vars, exp.AssignedValue)

			fieldType := structType.Fields[index]
			instance.UpdateField(block, index, x.Load(block), fieldType)
		}

		t.print(block, "instance - %s", instance.Name)
		fmt.Println("instance --- ", instance)
		return instance

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar := t.processExpression(block, vars, ex.Member)
		if baseVar == nil {
			panic(fmt.Sprintf("nil base in member expression: %v %v", ex.Member, vars))
		}

		// Base must be a class instance
		cls, ok := baseVar.(*Class)
		if !ok {
			panic(fmt.Sprintf("member access base is not a class instance, got %T", baseVar))
		}

		// Get metadata for base class
		classMeta, ok := t.classes[cls.Name]
		if !ok {
			panic(fmt.Sprintf("unknown class metadata: %s", cls.Name))
		}

		// Compute field name in identifier map
		fieldID := t.identifierBuilder.Attach(cls.Name, ex.Property)
		idx, ok := classMeta.varIndexMap[fieldID]
		if !ok {
			panic(fmt.Sprintf("unknown field %s on class %s", ex.Property, cls.Name))
		}

		// Get field type from struct UDT
		st, ok := classMeta.udt.(*types.StructType)
		if !ok {
			panic(fmt.Sprintf("class %s does not have a struct UDT", cls.Name))
		}
		fieldType := st.Fields[idx]

		// Get pointer to the field
		fieldPtr := cls.FieldPtr(block, idx)

		// Determine the class name if the field is a struct
		getClassName := func(tt types.Type) string {
			for cname, meta := range t.classes {
				if meta.udt == tt {
					return cname
				}
			}
			for cname, meta := range t.typeHandler.udts {
				if meta.udt == tt {
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
				return &Boolean{NativeType: types.I1, Value: fieldPtr}
			case 8:
				return &Int8{NativeType: types.I8, Value: fieldPtr}
			case 16:
				return &Int16{NativeType: types.I16, Value: fieldPtr}
			case 32:
				return &Int32{NativeType: types.I32, Value: fieldPtr}
			case 64:
				return &Int64{NativeType: types.I64, Value: fieldPtr}
			default:
				panic(fmt.Sprintf("unsupported int size %d", ft.BitSize))
			}

		case *types.FloatType:
			switch ft.Kind {
			case types.FloatKindHalf:
				return &Float16{NativeType: types.Half, Value: fieldPtr}
			case types.FloatKindFloat:
				return &Float32{NativeType: types.Float, Value: fieldPtr}
			case types.FloatKindDouble:
				return &Float64{NativeType: types.Double, Value: fieldPtr}
			default:
				panic(fmt.Sprintf("unsupported float kind %v", ft.Kind))
			}

		case *types.StructType:
			return &Class{
				Name: getClassName(fieldType),
				UDT:  fieldType,
				Ptr:  fieldPtr,
			}

		default:
			panic(fmt.Sprintf("unsupported field type %T in member expression", fieldType))
		}

	case ast.CallExpression:
		return t.handleCallExpression(block, vars, ex)

	case ast.BinaryExpression:
		left := t.processExpression(block, vars, ex.Left)
		right := t.processExpression(block, vars, ex.Right)
		if left == nil || right == nil {
			panic("nil operand in binary expression")
		}
		lv := left.Load(block)
		rv := right.Load(block)

		f := &Float64{}
		lvf := f.Cast(block, lv)
		rvf := f.Cast(block, rv)

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
		return &Float64{NativeType: types.Double, Value: ptr}
	}

	return nil
}

func (t *LLVM) handleCallExpression(block *ir.Block, vars map[string]Var, ex ast.CallExpression) Var {
	switch m := ex.Method.(type) {
	case ast.SymbolExpression:
		panic("method call should be on instance")

	case ast.MemberExpression:
		// Evaluate the base expression
		baseVar := t.processExpression(block, vars, m.Member)
		if baseVar == nil {
			panic("handleCallExpression: nil baseVar for member expression")
		}

		cls, ok := baseVar.(*Class)
		if !ok {
			panic(fmt.Sprintf("handleCallExpression: member access base is not Class (got %T)", baseVar))
		}
		if cls == nil || cls.Ptr == nil {
			panic(fmt.Sprintf("handleCallExpression: class or class.Ptr is nil for class %v", cls))
		}

		classMeta := t.classes[cls.Name]
		if classMeta == nil {
			panic("handleCallExpression: unknown class metadata: " + cls.Name)
		}

		methodKey := t.identifierBuilder.Attach(cls.Name, m.Property)
		fn, ok := classMeta.methods[methodKey]
		if !ok || fn == nil {
			panic(fmt.Sprintf("handleCallExpression: unknown method %s on class %s", m.Property, cls.Name))
		}

		// Build args for the user parameters (do not append `this` yet)
		args := make([]value.Value, 0, len(ex.Arguments)+1)
		for i, argExp := range ex.Arguments {
			v := t.processExpression(block, vars, argExp)
			if v == nil {
				panic(fmt.Sprintf("handleCallExpression: nil arg %d for %s.%s", i, cls.Name, m.Property))
			}
			raw := v.Load(block)
			if raw == nil {
				panic(fmt.Sprintf("handleCallExpression: loaded nil arg %d for %s.%s", i, cls.Name, m.Property))
			}

			// If the callee expects a certain param type, cast to it
			expected := fn.Sig.Params[i]
			raw = t.typeHandler.CastToType(block, expected.Name(), raw)
			if raw == nil {
				panic(fmt.Sprintf("handleCallExpression: CastToType returned nil for arg %d -> %s", i, expected.String()))
			}
			args = append(args, raw)
		}

		// Pass `this` as a pointer-to-struct (Slot returns pointer)
		thisPtr := cls.Slot()
		if thisPtr == nil {
			panic(fmt.Sprintf("handleCallExpression: this pointer is nil for %s", cls.Name))
		}

		// Check function expected param count: we declared 'this' last when creating fn,
		// adjust order according to how the function was declared.
		args = append(args, thisPtr)

		// ensure function is non-nil
		if fn == nil {
			panic(fmt.Sprintf("handleCallExpression: function pointer nil for %s.%s", cls.Name, m.Property))
		}

		// Now call
		ret := block.NewCall(fn, args...)
		if fn.Sig.RetType == types.Void {
			return nil
		}

		// allocate slot for return, store and wrap
		if ret == nil {
			panic(fmt.Sprintf("handleCallExpression: call returned nil for %s.%s", cls.Name, m.Property))
		}
		slot := block.NewAlloca(fn.Sig.RetType)
		if slot == nil {
			panic("handleCallExpression: failed to alloca for return")
		}
		block.NewStore(ret, slot)

		return t.wrapReturn(slot, fn.Sig.RetType, cls.Name)
	}
	return nil
}

func (t *LLVM) wrapReturn(slot *ir.InstAlloca, rt types.Type, className string) Var {
	switch v := rt.(type) {
	case *types.IntType:
		switch v.BitSize {
		case 1:
			return &Boolean{NativeType: types.I1, Value: slot}
		case 8:
			return &Int8{NativeType: types.I8, Value: slot}
		case 16:
			return &Int16{NativeType: types.I16, Value: slot}
		case 32:
			return &Int32{NativeType: types.I32, Value: slot}
		case 64:
			return &Int64{NativeType: types.I64, Value: slot}
		}
	case *types.FloatType:
		switch v.Kind {
		case types.FloatKindFloat:
			return &Float32{NativeType: types.Float, Value: slot}
		case types.FloatKindDouble:
			return &Float64{NativeType: types.Double, Value: slot}
		}
	case *types.StructType:
		return &Class{Name: className, UDT: v, Ptr: slot}
	}
	// fallback
	return &Class{Name: className, UDT: rt, Ptr: slot}
}
