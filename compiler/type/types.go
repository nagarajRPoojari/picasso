package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorutils "github.com/nagarajRPoojari/x-lang/compiler/error"
	rterr "github.com/nagarajRPoojari/x-lang/compiler/libs/private/runtime"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/boolean"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/floats"
	"github.com/nagarajRPoojari/x-lang/compiler/type/primitives/ints"
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

func NewTypeHandler() *TypeHandler {
	return &TypeHandler{Udts: make(map[string]*MetaClass)}
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
			return &boolean.Boolean{
				NativeType: types.I1,
				Value:      ptr,
				GoVal:      ci.X.Sign() != 0, // true if nonzero
			}
		}
		return &boolean.Boolean{NativeType: types.I1, Value: ptr, GoVal: false}

	case INT8, "i8":
		if init == nil {
			init = constant.NewInt(types.I8, 0)
		}
		ptr := block.NewAlloca(types.I8)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &ints.Int8{
				NativeType: types.I8,
				Value:      ptr,
				GoVal:      int8(ci.X.Int64()),
			}
		}
		return &ints.Int8{NativeType: types.I8, Value: ptr, GoVal: 0}

	case INT16, "i16":
		if init == nil {
			init = constant.NewInt(types.I16, 0)
		}
		ptr := block.NewAlloca(types.I16)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &ints.Int16{
				NativeType: types.I16,
				Value:      ptr,
				GoVal:      int16(ci.X.Int64()),
			}
		}
		return &ints.Int16{NativeType: types.I16, Value: ptr, GoVal: 0}

	case INT32, "i32":
		if init == nil {
			init = constant.NewInt(types.I32, 0)
		}
		ptr := block.NewAlloca(types.I32)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &ints.Int32{
				NativeType: types.I32,
				Value:      ptr,
				GoVal:      int32(ci.X.Int64()),
			}
		}
		return &ints.Int32{NativeType: types.I32, Value: ptr, GoVal: 0}

	case INT64, INT, "i64":
		if init == nil {
			init = constant.NewInt(types.I64, 0)
		}
		ptr := block.NewAlloca(types.I64)
		block.NewStore(init, ptr)

		if ci, ok := init.(*constant.Int); ok {
			return &ints.Int64{
				NativeType: types.I64,
				Value:      ptr,
				GoVal:      ci.X.Int64(),
			}
		}
		return &ints.Int64{NativeType: types.I64, Value: ptr, GoVal: 0}

	case FLOAT16, "half":
		if init == nil {
			init = constant.NewFloat(types.Half, 0.0)
		}
		ptr := block.NewAlloca(types.Half)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &floats.Float16{
				NativeType: types.Half,
				Value:      ptr,
				GoVal:      float32(f),
			}
		}
		return &floats.Float16{NativeType: types.Float, Value: ptr, GoVal: 0}

	case FLOAT32, "float":
		if init == nil {
			init = constant.NewFloat(types.Float, 0.0)
		}
		ptr := block.NewAlloca(types.Float)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &floats.Float32{
				NativeType: types.Float,
				Value:      ptr,
				GoVal:      float32(f),
			}
		}
		return &floats.Float32{NativeType: types.Float, Value: ptr, GoVal: 0}

	case FLOAT64, DOUBLE:
		if init == nil {
			init = constant.NewFloat(types.Double, 0.0)
		}
		ptr := block.NewAlloca(types.Double)
		block.NewStore(init, ptr)

		if cf, ok := init.(*constant.Float); ok {
			f, _ := cf.X.Float64()
			return &floats.Float64{
				NativeType: types.Double,
				Value:      ptr,
				GoVal:      f,
			}
		}
		return &floats.Float64{NativeType: types.Double, Value: ptr, GoVal: 0}
	case STRING, "i8*":
		if init == nil {
			init = constant.NewNull(types.I8Ptr)
		}
		return NewString(block, init)
	case NULL, VOID:
		return NewNullVar(types.NewPointer(init.Type()))
	}

	if udt, ok := t.Udts[string(_type)]; ok {
		if init == nil {
			init = constant.NewZeroInitializer(udt.UDT)
		}
		c := NewClass(
			block, string(_type), udt.UDT,
		)
		c.Update(block, init)
		return c
	}

	errorutils.Abort(errorutils.TypeError, errorutils.InvalidNativeType, _type)
	return nil
}

// GetLLVMType accepts native Type & returns llvm compatible types.Type
func (t *TypeHandler) GetLLVMType(_type Type) types.Type {
	switch _type {
	case NULL, VOID:
		return types.Void
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
		return types.I8Ptr
	case ARRAY:
		// return
	}

	// Check if already registered
	if k, ok := t.Udts[string(_type)]; ok {
		return k.UDT
	}

	errorutils.Abort(errorutils.TypeError, errorutils.InvalidLLVMType, _type)
	return nil
}

// ImplicitTypeCast takes a target type name (e.g. "float64", "int8")
// and a value, and emits the appropriate cast instruction in `block`.
func (t *TypeHandler) ImplicitTypeCast(block *ir.Block, target string, v value.Value) (value.Value, *ir.Block) {
	switch target {
	case "boolean", "bool", "i1":
		return t.ImplicitIntCast(block, v, types.I1)
	case "int8", "i8":
		return t.ImplicitIntCast(block, v, types.I8)
	case "int16", "i16":
		return t.ImplicitIntCast(block, v, types.I16)
	case "int32", "i32":
		return t.ImplicitIntCast(block, v, types.I32)
	case "int", "int64", "i64":
		return t.ImplicitIntCast(block, v, types.I64)

	case "float16", "half":
		return t.ImplicitFloatCast(block, v, types.Half)
	case "float32", "float":
		return t.ImplicitFloatCast(block, v, types.Float)
	case "float64", "double":
		return t.ImplicitFloatCast(block, v, types.Double)
	case "string", "i8*":
		switch v.Type().(type) {
		case *types.PointerType:
			return v, block
		default:
			errorsx.PanicCompilationError(fmt.Sprintf(
				"cannot cast %s to string", v.Type().String(),
			))
		}
	case "void":
		return nil, block
	}

	if k, ok := t.Udts[target]; ok {
		return ensureType(block, v, k.UDT), block
	}
	errorutils.Abort(errorutils.TypeError, errorutils.InvalidTargetType, target)
	return nil, block
}

// NOTE: assumes floatMax/floatMin and intMax/intMin maps exist as you already have.

func (t *TypeHandler) catchIntToIntDownCast(block *ir.Block, v value.Value, dst *types.IntType) (value.Value, *ir.Block) {
	b := block

	// boolean target (i1): non-zero -> true
	if dst.BitSize == 1 {
		cmp := b.NewICmp(enum.IPredNE, v, constant.NewInt(v.Type().(*types.IntType), 0))
		return cmp, b
	}

	abort := b.Parent.NewBlock("")
	safe := b.Parent.NewBlock("")

	maxVal := constant.NewInt(dst, intMax[dst])
	minVal := constant.NewInt(dst, intMin[dst])
	overflowMax := b.NewICmp(enum.IPredSGT, v, maxVal)
	overflowMin := b.NewICmp(enum.IPredSLT, v, minVal)
	overflow := b.NewOr(overflowMax, overflowMin)

	b.NewCondBr(overflow, abort, safe)

	rterr.Instance.RaiseRTError(abort, "runtime overflow in int downcast\n")
	abort.NewUnreachable()

	vTrunc := safe.NewTrunc(v, dst)
	return vTrunc, safe
}

func (t *TypeHandler) catchFloatToIntDownCast(block *ir.Block, v value.Value, dst *types.IntType) (value.Value, *ir.Block) {
	b := block

	abort := b.Parent.NewBlock("")
	safe := b.Parent.NewBlock("")

	// Promote the float to double for robust comparison (avoids half literal issues)
	var vAsDouble value.Value
	if _, ok := v.Type().(*types.FloatType); ok && v.Type().(*types.FloatType).Kind != types.FloatKindDouble {
		vAsDouble = b.NewFPExt(v, types.Double)
	} else {
		vAsDouble = v
	}

	// Compare against integer bounds expressed as double constants
	minValD := constant.NewFloat(types.Double, float64(intMin[dst]))
	maxValD := constant.NewFloat(types.Double, float64(intMax[dst]))

	overflowMax := b.NewFCmp(enum.FPredOGT, vAsDouble, maxValD)
	overflowMin := b.NewFCmp(enum.FPredOLT, vAsDouble, minValD)
	overflow := b.NewOr(overflowMax, overflowMin)

	b.NewCondBr(overflow, abort, safe)

	rterr.Instance.RaiseRTError(abort, "runtime overflow in float → int downcast\n")
	abort.NewUnreachable()

	// In safe block: do the FP->SI conversion. Use vAsDouble (double->int) which is fine.
	// If you prefer converting from original width, you can do safe.NewFPToSI(v, dst) too;
	// FPToSI from double is fine and avoids repeated casting complexities.
	res := safe.NewFPToSI(vAsDouble, dst)
	return res, safe
}

func (t *TypeHandler) ImplicitIntCast(block *ir.Block, v value.Value, dst *types.IntType) (value.Value, *ir.Block) {
	b := block
	src, ok := v.Type().(*types.IntType)
	if !ok {
		// v is not an int; allow float -> int/boolean conversions
		if _, ok := v.Type().(*types.FloatType); ok {
			// float -> boolean
			if dst == types.I1 {
				zero := constant.NewFloat(v.Type().(*types.FloatType), 0.0)
				cond := b.NewFCmp(enum.FPredONE, v, zero) // v != 0.0
				return cond, b
			}

			// float -> integer (with overflow checks)
			return t.catchFloatToIntDownCast(block, v, dst)
		}
		panic("cannot intCast from " + v.Type().String())
	}

	if src.BitSize == 1 {
		if dst.BitSize == 1 {
			return v, b
		}
		// Boolean widen → ZExt
		return b.NewZExt(v, dst), b
	}

	// v is integer
	if src.BitSize > dst.BitSize {
		// downcast integer -> integer (with overflow checks)
		return t.catchIntToIntDownCast(block, v, dst)
	}
	if src.BitSize < dst.BitSize {
		// extend smaller integer to larger
		return b.NewSExt(v, dst), b
	}
	return v, b
}

func (t *TypeHandler) catchFloatToFloatDowncast(block *ir.Block, v value.Value, src *types.FloatType, dst *types.FloatType) (value.Value, *ir.Block) {
	b := block

	// If same type, no cast needed
	if src.Kind == dst.Kind {
		return v, b
	}

	// Upcast (smaller -> larger)
	if floatRank(src.Kind) < floatRank(dst.Kind) {
		return b.NewFPExt(v, dst), b
	}

	// Downcast (larger -> smaller)
	abort := b.Parent.NewBlock("")
	safe := b.Parent.NewBlock("")

	// Promote v to double for comparisons (avoids half literal problems)
	var vAsDouble value.Value
	if src.Kind != types.FloatKindDouble {
		vAsDouble = b.NewFPExt(v, types.Double)
	} else {
		vAsDouble = v
	}

	// Use floatMax/floatMin maps (float64 values) as double constants
	maxD := constant.NewFloat(types.Double, floatMax[dst])
	minD := constant.NewFloat(types.Double, floatMin[dst])

	overflowMax := b.NewFCmp(enum.FPredOGT, vAsDouble, maxD)
	overflowMin := b.NewFCmp(enum.FPredOLT, vAsDouble, minD)
	overflow := b.NewOr(overflowMax, overflowMin)

	b.NewCondBr(overflow, abort, safe)

	rterr.Instance.RaiseRTError(abort, "runtime overflow in float demotion")
	abort.NewUnreachable()

	// Safe block: actually truncate original value to dst type
	vTrunc := safe.NewFPTrunc(v, dst)
	return vTrunc, safe
}

func (t *TypeHandler) catchIntToFloatDowncast(block *ir.Block, v value.Value, dst *types.FloatType) (value.Value, *ir.Block) {
	b := block

	abort := b.Parent.NewBlock("")
	safe := b.Parent.NewBlock("")

	// Convert integer to double for comparison (avoids half literal issues)
	vAsDouble := b.NewSIToFP(v, types.Double)

	maxD := constant.NewFloat(types.Double, floatMax[dst])
	minD := constant.NewFloat(types.Double, floatMin[dst])

	overflowMax := b.NewFCmp(enum.FPredOGT, vAsDouble, maxD)
	overflowMin := b.NewFCmp(enum.FPredOLT, vAsDouble, minD)
	overflow := b.NewOr(overflowMax, overflowMin)

	b.NewCondBr(overflow, abort, safe)

	rterr.Instance.RaiseRTError(abort, "runtime overflow converting int → float")
	abort.NewUnreachable()

	// Safe block: return converted float in requested dst width
	res := safe.NewSIToFP(v, dst)
	return res, safe
}

func (t *TypeHandler) ImplicitFloatCast(block *ir.Block, v value.Value, dst *types.FloatType) (value.Value, *ir.Block) {
	switch src := v.Type().(type) {
	case *types.FloatType:
		return t.catchFloatToFloatDowncast(block, v, src, dst)

	case *types.IntType:
		// int -> float: special-case i1 -> treat as 0/1
		if src.BitSize == 1 {
			intVal := block.NewZExt(v, types.I8) // 0 or 1 -> i8 then convert
			floatVal := block.NewSIToFP(intVal, dst)
			return floatVal, block
		}
		return t.catchIntToFloatDowncast(block, v, dst)

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
