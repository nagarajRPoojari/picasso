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
	nonStaticVarsIndices map[string]int
	nonStaticVars        []constant.Constant
	nonStaticVarsAST     map[string]*ast.VariableDeclarationStatement

	staticMethods    map[string]*ir.Func
	nonStaticMethods map[string]*ir.Func

	udt types.Type
}

func NewMetaClass() *MetaClass {
	return &MetaClass{
		nonStaticVars:        make([]constant.Constant, 0),
		nonStaticVarsIndices: make(map[string]int),
		staticMethods:        make(map[string]*ir.Func),
		nonStaticMethods:     make(map[string]*ir.Func),
		nonStaticVarsAST:     make(map[string]*ast.VariableDeclarationStatement),
	}
}

type LLVM struct {
	module *ir.Module

	typeHandler       *TypeHandler
	identifierBuilder *IdentifierBuilder

	vars    map[string]Var
	methods map[string]*ir.Func
	classes map[string]*MetaClass

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
	f, err := os.Create("output.ll")
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
			t.classes[st.Name] = NewMetaClass()
			t.initClass(st)
		}
	}

	for _, stI := range tree.Body {
		switch st := stI.(type) {
		case ast.ClassDeclarationStatement:
			t.declareStaticFields(st)
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
				t.defineMain(&st)
			}
		}
	}
}

func (t *LLVM) initClass(class ast.ClassDeclarationStatement) {
	i := 0
	tps := make([]types.Type, 0)
	for _, stI := range class.Body {
		switch st := stI.(type) {
		// @todo: casting constant.Constant type
		case ast.VariableDeclarationStatement:
			if !st.IsStatic {
				name := t.identifierBuilder.Attach(class.Name, st.Identifier)
				t.classes[class.Name].nonStaticVarsIndices[name] = i
				i++

				v := t.processStaticExpression(st.AssignedValue)
				t.classes[class.Name].nonStaticVars = append(t.classes[class.Name].nonStaticVars, v)
				tps = append(tps, v.Type())
			}
		}
	}

	udt := types.NewStruct(tps...)
	t.module.NewTypeDef(class.Name, udt)
	t.classes[class.Name].udt = udt
}

func (t *LLVM) declareStaticFields(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			if st.IsStatic {
				if st.AssignedValue == nil {
					panic("static varables should be initialized")
				}
				name := t.identifierBuilder.Attach(class.Name, st.Identifier)
				t.module.NewGlobalDef(name, t.processStaticExpression(st.AssignedValue))
			}
		case ast.FunctionDeclarationStatement:
			params := make([]*ir.Param, 0)
			for _, p := range st.Parameters {
				params = append(params, ir.NewParam(p.Name, t.typeHandler.GetLLVMType(Type(p.Type.Get()))))
			}
			if !st.IsStatic {
				params = append(params, ir.NewParam("this", t.classes[class.Name].udt))
			}
			name := t.identifierBuilder.Attach(class.Name, st.Name)
			b := t.typeHandler.GetLLVMType(Type(st.ReturnType.Get()))
			f := t.module.NewFunc(name, b, params...)
			t.methods[name] = f
			if st.IsStatic {
				t.classes[class.Name].staticMethods[name] = f
			} else {
				t.classes[class.Name].nonStaticMethods[name] = f
			}
		}
	}
}

func (t *LLVM) processStaticExpression(expI ast.Expression) constant.Constant {
	switch exp := expI.(type) {
	case ast.NumberExpression:
		return constant.NewFloat(types.Double, float64(exp.Value))

	case ast.BinaryExpression:
		left := t.processStaticExpression(exp.Left)
		right := t.processStaticExpression(exp.Right)

		if left == nil || right == nil {
			return nil
		}

		lf, lok := left.(*constant.Float)
		rf, rok := right.(*constant.Float)
		if !lok || !rok {
			return nil
		}

		switch exp.Operator.Value {
		case "+":
			return constant.NewFAdd(lf, rf)
		case "-":
			return constant.NewFSub(lf, rf)
		case "*":
			return constant.NewFMul(lf, rf)
		case "/":
			return constant.NewFDiv(lf, rf)
		default:
			return nil
		}
	}
	return nil
}

func (t *LLVM) defineClass(class ast.ClassDeclarationStatement) {
	for _, stI := range class.Body {
		switch st := stI.(type) {
		case ast.FunctionDeclarationStatement:
			t.defineFunc(class.Name, &st)
		}
	}
}

func (t *LLVM) defineMain(fn *ast.FunctionDeclarationStatement) {
	vars := make(map[string]Var, 0)
	f := t.methods[MAIN]
	entry := f.NewBlock(ENTRY)

	for _, stI := range fn.Body {
		switch st := stI.(type) {
		case ast.VariableDeclarationStatement:
			r := t.processExpression(entry, vars, st.AssignedValue)
			if _, ok := vars[st.Identifier]; ok {
				panic("variable already exists")
			}
			vars[st.Identifier] = r

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
				rhs := t.processStaticExpression(exp.AssignedValue)
				v.Update(entry, rhs)
			case ast.CallExpression:
				t.processExpression(entry, vars, exp)
			default:
				panic("invalid statement")
			}
		}
	}

	val := vars["z"].Load(entry)
	t.print(entry, "final == %f ", val)

	entry.NewRet(constant.NewInt(types.I32, 0))
}

func (t *LLVM) defineFunc(className string, fn *ast.FunctionDeclarationStatement) {
	vars := make(map[string]Var, 0)
	name := t.identifierBuilder.Attach(className, fn.Name)
	f := t.methods[name]
	entry := f.NewBlock("entry")

	for i, p := range f.Params {
		if i >= len(fn.Parameters) {
			this := f.Params[len(f.Params)-1]
			cls := t.classes[className]
			vars[this.LocalName] = NewClass(entry, className, cls.udt)
			break
		}
		paramType := Type(fn.Parameters[i].Type.Get())
		vars[p.LocalName] = t.typeHandler.GetPrimitiveVar(entry, paramType, p)
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
				v = t.typeHandler.GetPrimitiveVar(entry, Type(st.ExplicitType.Get()), casted)
			}

			vars[st.Identifier] = v

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
		// instantiate class and initialize fields from arguments
		meth := ex.Instantiation.Method.(ast.SymbolExpression)
		classMeta := t.classes[meth.Value]
		if classMeta == nil {
			panic(fmt.Sprintf("unknown class: %s", meth.Value))
		}
		instance := NewClass(block, meth.Value, classMeta.udt)
		for i, v := range t.classes[meth.Value].nonStaticVars {
			if st, ok := classMeta.udt.(*types.StructType); ok {
				fieldTy := st.Fields[i]
				argCast := t.typeHandler.CastToType(block, fieldTy.String(), v)
				instance.UpdateField(block, i, argCast)
			} else {
				instance.UpdateField(block, i, v)
			}
		}

		return instance

	case ast.MemberExpression:
		baseVar := t.processExpression(block, vars, ex.Member)
		if baseVar == nil {
			panic("nil base in member expression")
		}

		cls, ok := baseVar.(*Class)
		if !ok {
			panic("member access base is not a class instance")
		}

		classMeta := t.classes[cls.Name]
		if classMeta == nil {
			panic("unknown class metadata: " + cls.Name)
		}

		property := t.identifierBuilder.Attach(cls.Name, ex.Property)
		idx, ok := classMeta.nonStaticVarsIndices[property]
		if !ok {
			panic(fmt.Sprintf("unknown field %s on class %s", ex.Property, cls.Name))
		}

		// field type comes from struct UDT fields if UDT is struct
		var fieldType types.Type
		if st, ok := classMeta.udt.(*types.StructType); ok {
			fieldType = st.Fields[idx]
		} else {
			// fallback: assume i8* or something reasonable
			fieldType = types.I8
		}

		// Get pointer to field then load
		fieldPtr := cls.FieldPtr(block, idx)
		fieldVal := block.NewLoad(fieldType, fieldPtr)

		// Create a new alloca slot and store the field value there so we return a mutable Var.
		slot := block.NewAlloca(fieldType)
		block.NewStore(fieldVal, slot)

		// Wrap into appropriate Var implementation
		switch ft := fieldType.(type) {
		case *types.IntType:
			switch ft.BitSize {
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
			switch ft.Kind {
			case types.FloatKindFloat:
				return &Float32{NativeType: types.Float, Value: slot}
			case types.FloatKindDouble:
				return &Float64{NativeType: types.Double, Value: slot}
			default:
				return &Float64{NativeType: types.Double, Value: slot}
			}

		default:
			return &Class{UDT: fieldType, Value: slot, Name: cls.Name}
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
		// First evaluate the base
		baseVar := t.processExpression(block, vars, m.Member)
		if baseVar == nil {
			panic("nil base in member expression")
		}

		cls, ok := baseVar.(*Class)
		if !ok {
			panic(fmt.Sprintf("member access base is not a class instance (got %T)", baseVar))
		}

		classMeta := t.classes[cls.Name]
		if classMeta == nil {
			panic("unknown class metadata: " + cls.Name)
		}

		methodKey := t.identifierBuilder.Attach(cls.Name, m.Property)

		// Instance method
		if fn, ok := classMeta.nonStaticMethods[methodKey]; ok {
			args := make([]value.Value, 0, len(ex.Arguments)+1)

			// user args
			for i, argExp := range ex.Arguments {
				v := t.processExpression(block, vars, argExp)
				raw := v.Load(block)

				// Cast arg to expected param type
				if i < len(fn.Sig.Params) {
					expected := fn.Sig.Params[i]
					raw = t.typeHandler.CastToType(block, expected.String(), raw)
				}
				args = append(args, raw)
			}

			// "this" pointer goes last
			thisVal := cls.Load(block)
			args = append(args, thisVal)

			ret := block.NewCall(fn, args...)

			if fn.Sig.RetType == types.Void {
				return nil
			}

			slot := block.NewAlloca(fn.Sig.RetType)
			block.NewStore(ret, slot)

			return t.wrapReturn(slot, fn.Sig.RetType, cls.Name)
		}

		// Static method
		if fn, ok := classMeta.staticMethods[methodKey]; ok {
			args := make([]value.Value, 0, len(ex.Arguments))
			for i, argExp := range ex.Arguments {
				v := t.processExpression(block, vars, argExp)
				raw := v.Load(block)

				// Cast to expected param type
				if i < len(fn.Sig.Params) {
					expected := fn.Sig.Params[i]
					raw = t.typeHandler.CastToType(block, expected.String(), raw)
				}
				args = append(args, raw)
			}

			ret := block.NewCall(fn, args...)
			if fn.Sig.RetType == types.Void {
				return nil
			}

			slot := block.NewAlloca(fn.Sig.RetType)
			block.NewStore(ret, slot)

			return t.wrapReturn(slot, fn.Sig.RetType, cls.Name)
		}

		panic(fmt.Sprintf("unknown method %s on class %s", m.Property, cls.Name))
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
		return &Class{Name: className, UDT: v, Value: slot}
	}
	// fallback
	return &Class{Name: className, UDT: rt, Value: slot}
}
