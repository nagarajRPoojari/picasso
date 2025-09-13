package compiler

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
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
	NULL   Type = "null"
)

type TypeHandler struct {
	udts map[string]*Class
}

func NewTypeHandler() *TypeHandler {
	return &TypeHandler{
		udts: map[string]*Class{},
	}
}

func (t *TypeHandler) AddUDT(name string, c *Class) {
	t.udts[name] = c
}

func (t *TypeHandler) GetUDT(name string) {

}

func (t *TypeHandler) GetPrimitiveVar(block *ir.Block, _type Type, init value.Value) Var {
	switch _type {
	case BOOLEAN:
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

	case INT8:
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

	case INT16:
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

	case INT32:
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

	case INT64, INT:
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

	case FLOAT32:
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

	default:
		return nil
	}
}

func (t *TypeHandler) GetVarType(_type Type) types.Type {
	switch _type {
	case BOOLEAN:
		return types.I1
	case INT8:
		return types.I8
	case INT16:
		return types.I16
	case INT32:
		return types.I32
	case INT64, INT:
		return types.I64
	case FLOAT16:
		return types.Half
	case FLOAT32:
		return types.Float
	case FLOAT64, DOUBLE:
		return types.Double
	case STRING:
		return types.NewPointer(types.I8)
	default:
		return nil
	}
}

// helper: bool → int64
func btoi(b bool) int64 {
	if b {
		return 1
	}
	return 0
}
