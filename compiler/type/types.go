package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

type Type string

const (
	BOOLEAN Type = "boolean"

	INT   Type = "int"
	INT8  Type = "int8"
	INT16 Type = "int16"
	INT32 Type = "int32"
	INT64 Type = "int64"

	FLOAT16 Type = "float16"
	FLOAT32 Type = "float32"
	FLOAT64 Type = "float64"

	DOUBLE Type = "double"

	STRING Type = "string"

	ARRAY Type = "array"

	NULL Type = "null"
	VOID Type = "void"
)

type TypeHandler struct {
	Udts   map[string]*MetaClass
	module *ir.Module

	langAlloc *ir.Func
}

func NewTypeHandler(l *ir.Func) *TypeHandler {
	return &TypeHandler{Udts: make(map[string]*MetaClass), langAlloc: l}
}

func (t *TypeHandler) Register(name string, meta *MetaClass) {
	t.Udts[name] = meta
}

func (t *TypeHandler) BuildVar(block *ir.Block, _type Type, init value.Value) Var {
	switch _type {
	case BOOLEAN, "i1":
		if init == nil {
			init = constant.NewInt(types.I1, 0)
		}
		ptr := block.NewAlloca(types.I1)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &Boolean{
				NativeType: types.I1,
				Value:      ptr,
				GoVal:      ci.X.Sign() != 0, // true if nonzero
			}
		}
		return &Boolean{NativeType: types.I1, Value: ptr, GoVal: false}

	case INT8, "i8":
		if init == nil {
			init = constant.NewInt(types.I8, 0)
		}
		ptr := block.NewAlloca(types.I8)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &Int8{
				NativeType: types.I8,
				Value:      ptr,
				GoVal:      int8(ci.X.Int64()),
			}
		}
		return &Int8{NativeType: types.I8, Value: ptr, GoVal: 0}

	case INT16, "i16":
		if init == nil {
			init = constant.NewInt(types.I16, 0)
		}
		ptr := block.NewAlloca(types.I16)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &Int16{
				NativeType: types.I16,
				Value:      ptr,
				GoVal:      int16(ci.X.Int64()),
			}
		}
		return &Int16{NativeType: types.I16, Value: ptr, GoVal: 0}

	case INT32, "i32":
		if init == nil {
			init = constant.NewInt(types.I32, 0)
		}
		ptr := block.NewAlloca(types.I32)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &Int32{
				NativeType: types.I32,
				Value:      ptr,
				GoVal:      int32(ci.X.Int64()),
			}
		}
		return &Int32{NativeType: types.I32, Value: ptr, GoVal: 0}

	case INT64, INT, "i64":
		if init == nil {
			init = constant.NewInt(types.I64, 0)
		}
		ptr := block.NewAlloca(types.I64)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &Int64{
				NativeType: types.I64,
				Value:      ptr,
				GoVal:      ci.X.Int64(),
			}
		}
		return &Int64{NativeType: types.I64, Value: ptr, GoVal: 0}

	case FLOAT16, "half":
		if init == nil {
			init = constant.NewFloat(types.Half, 0.0)
		}
		ptr := block.NewAlloca(types.Half)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &Float16{
				NativeType: types.Half,
				Value:      ptr,
				GoVal:      float32(f),
			}
		}
		return &Float16{NativeType: types.Float, Value: ptr, GoVal: 0}

	case FLOAT32, "float":
		if init == nil {
			init = constant.NewFloat(types.Float, 0.0)
		}
		ptr := block.NewAlloca(types.Float)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &Float32{
				NativeType: types.Float,
				Value:      ptr,
				GoVal:      float32(f),
			}
		}
		return &Float32{NativeType: types.Float, Value: ptr, GoVal: 0}

	case FLOAT64, DOUBLE:
		if init == nil {
			init = constant.NewFloat(types.Double, 0.0)
		}
		ptr := block.NewAlloca(types.Double)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &Float64{
				NativeType: types.Double,
				Value:      ptr,
				GoVal:      f,
			}
		}
		return &Float64{NativeType: types.Double, Value: ptr, GoVal: 0}
	case STRING:
		if init == nil {
			init = constant.NewNull(types.I8Ptr)
		}
		return NewString(block, init)
	case NULL, VOID:
		return NewNullVar(types.NewPointer(init.Type()))
	}

	if udt, ok := t.Udts[string(_type)]; ok {
		if init == nil {
			constant.NewZeroInitializer(udt.UDT)
		}
		c := NewClass(
			block, string(_type), udt.UDT, t.langAlloc,
		)
		c.Update(block, init)
		return c
	}

	errorsx.PanicCompilationError(fmt.Sprintf("invalid primitive type: %s , %s", _type, t.Udts))
	return nil
}

// GetLLVMType accepts native Type & returns llvm compatible types.Type
func (t *TypeHandler) GetLLVMType(_type Type) types.Type {
	switch _type {
	case NULL, VOID:
		return types.NewPointer(types.NewStruct())
	case BOOLEAN, "i1":
		return types.I1
	case INT8, "i8":
		return types.I8
	case INT16, "i16":
		return types.I16
	case INT32, "132":
		return types.I32
	case INT64, INT, "i64":
		return types.I64
	case FLOAT16, "half":
		return types.Half
	case FLOAT32, "float":
		return types.Float
	case FLOAT64, DOUBLE:
		return types.Double
	case STRING:
		return types.NewPointer(types.I8Ptr)
	case ARRAY:
		// return
	}

	// Check if already registered
	if k, ok := t.Udts[string(_type)]; ok {
		return k.UDT
	}

	errorsx.PanicCompilationError((fmt.Sprintf("invalid LLVM type: %s", _type)))
	return nil
}

// CastToType takes a target type name (e.g. "float64", "int8")
// and a value, and emits the appropriate cast instruction in `block`.
func (t *TypeHandler) CastToType(block *ir.Block, target string, v value.Value) value.Value {
	switch target {
	case "boolean", "bool", "i1":
		if v.Type().Equal(types.I1) {
			return v
		}
		switch v.Type().(type) {
		case *types.IntType:
			zero := constant.NewInt(v.Type().(*types.IntType), 0)
			return block.NewICmp(enum.IPredNE, v, zero)
		case *types.FloatType:
			zero := constant.NewFloat(v.Type().(*types.FloatType), 0.0)
			return block.NewFCmp(enum.FPredONE, v, zero)
		default:
			panic("cannot cast to boolean from type " + v.Type().String())
		}

	case "int8", "i8":
		return t.intCast(block, v, types.I8)
	case "int16", "i16":
		return t.intCast(block, v, types.I16)
	case "int32", "i32":
		return t.intCast(block, v, types.I32)
	case "int", "int64", "i64":
		return t.intCast(block, v, types.I64)

	case "float16", "half":
		return t.floatCast(block, v, types.Half)
	case "float32", "float":
		return t.floatCast(block, v, types.Float)
	case "float64", "double":
		return t.floatCast(block, v, types.Double)
	case "string":
		switch v.Type().(type) {
		case *types.PointerType:
			return v
		default:
			errorsx.PanicCompilationError(fmt.Sprintf(
				"cannot cast %s to string", v.Type().String(),
			))
		}
	}

	if k, ok := t.Udts[target]; ok {
		return ensureType(block, v, k.UDT)
	}
	errorsx.PanicCompilationError(fmt.Sprintf("unexpected target type: %s", target))
	return nil
}

func (t *TypeHandler) intCast(block *ir.Block, v value.Value, dst *types.IntType) value.Value {
	src, ok := v.Type().(*types.IntType)
	if !ok {
		// int ← float
		if _, ok := v.Type().(*types.FloatType); ok {
			return block.NewFPToSI(v, dst)
		}
		panic("cannot intCast from " + v.Type().String())
	}
	if src.BitSize > dst.BitSize {
		return block.NewTrunc(v, dst)
	} else if src.BitSize < dst.BitSize {
		return block.NewSExt(v, dst)
	}
	return v
}

func (t *TypeHandler) floatCast(block *ir.Block, v value.Value, dst *types.FloatType) value.Value {
	switch src := v.Type().(type) {
	case *types.FloatType:
		if src.Kind == dst.Kind {
			return v
		}
		// Promote/demote based on known float kinds
		if floatRank(src.Kind) < floatRank(dst.Kind) {
			return block.NewFPExt(v, dst) // promote
		}
		return block.NewFPTrunc(v, dst) // demote

	case *types.IntType:
		return block.NewSIToFP(v, dst) // signed int to float

	default:
		panic("cannot floatCast from " + v.Type().String())
	}
}

func floatRank(k types.FloatKind) int {
	switch k {
	case types.FloatKindHalf:
		return 16
	case types.FloatKindFloat:
		return 32
	case types.FloatKindDouble:
		return 64
	case types.FloatKindX86_FP80:
		return 80
	case types.FloatKindFP128, types.FloatKindPPC_FP128:
		return 128
	default:
		return 0
	}
}
