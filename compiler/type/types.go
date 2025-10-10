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

type VarBlock struct {
	*ir.Block
}

type BlockHolder struct {
	V VarBlock
	N *ir.Block
}

type Type struct {
	T string
	U string
}

func NewType(T string, U ...string) Type {
	if len(U) > 0 {
		return Type{T: T, U: U[0]}
	}
	return Type{T: T}
}

const (
	BOOLEAN = "boolean"

	INT   = "int"
	INT8  = "int8"
	INT16 = "int16"
	INT32 = "int32"
	INT64 = "int64"

	FLOAT16 = "float16"
	FLOAT32 = "float32"
	FLOAT64 = "float64"

	DOUBLE = "double"

	STRING = "string"

	ARRAY = "array"

	NULL = "null"
	VOID = "void"
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

func (t *TypeHandler) Exists(tp string) bool {
	switch tp {
	case NULL, VOID, BOOLEAN, "i1", INT8, "i8", INT16, "i16", INT32, "132", INT64, INT, "i64", FLOAT16, "half", FLOAT32, "float", FLOAT64, DOUBLE, STRING:
		return true
	}

	// Check if already registered
	if _, ok := t.Udts[tp]; ok {
		return true
	}

	return false
}

// BuildVar creates and initializes a new variable of the given type in the
// specified LLVM IR block. It allocates storage, applies an optional
// initializer, and returns a Var wrapper that provides runtime access.
//
// Parameters:
//
//	block — the LLVM IR basic block where the variable is allocated.
//	_type — the high-level type identifier (primitive, string, void, or Class type).
//	init  — optional initializer value; if nil, a default zero value is used.
//
// Returns:
//
//	Var — a wrapper around the allocated variable.
//
// Note:
//   - class must be registered with TypeHandler before building var.
func (t *TypeHandler) BuildVar(bh BlockHolder, _type Type, init value.Value) Var {
	switch _type.T {
	case BOOLEAN, "i1":
		if init == nil {
			init = constant.NewInt(types.I1, 0)
		}
		ptr := bh.V.NewAlloca(types.I1)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.I8)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.I16)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.I32)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.I64)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.Half)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.Float)
		bh.N.NewStore(init, ptr)

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
		ptr := bh.V.NewAlloca(types.Double)
		bh.N.NewStore(init, ptr)

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
		return NewString(bh, init)
	case NULL, VOID:
		return NewNullVar(types.NewPointer(init.Type()))
	case ARRAY:
		if _type.U == "" {
			errorutils.Abort(errorutils.InternalError, errorutils.InternalError, "sub type should be provided for array type")
		}

		ele := t.GetLLVMType(_type.U)
		if init == nil {
			pt := types.NewPointer(ARRAYSTRUCT)
			init = constant.NewNull(pt)
		}
		return &Array{
			Ptr:       init,
			ArrayType: ARRAYSTRUCT,
			ElemType:  ele,
		}
	}

	if udt, ok := t.Udts[string(_type.T)]; ok {
		if init == nil {
			init = constant.NewZeroInitializer(udt.UDT)
		}
		c := NewClass(
			bh.V, string(_type.T), udt.UDT,
		)
		c.Update(bh.N, init)
		return c
	}

	errorutils.Abort(errorutils.TypeError, errorutils.InvalidNativeType, _type)
	return nil
}

// GetLLVMType maps a high-level type identifier to its corresponding LLVM type.
//
// Parameters:
//
//	_type — the high-level type identifier (e.g., INT32, FLOAT64, STRING, UDT).
//
// Returns:
//
//	types.Type — the LLVM-compatible type that matches the given high-level type.
//
// Special cases:
//   - NULL, VOID → types.Void
//   - Boolean → types.I1
//   - Integers → LLVM integer types (I8, I16, I32, I64)
//   - Floats → LLVM floating-point types (Half, Float, Double)
//   - String → i8 pointer
//   - UDTs → resolved from the registered type table
//
// If the type is unknown or unsupported, the function aborts with a type error.
func (t *TypeHandler) GetLLVMType(_type string) types.Type {
	if _type == "" {
		return types.Void
	}

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
		s := types.NewStruct(
			types.I64,                   // length
			types.NewPointer(types.I8),  // data
			types.NewPointer(types.I64), // shape (i64*)
			types.I64,                   // rank
		)

		s.SetName("array")
		return types.NewPointer(s)
	}

	// Check if already registered
	if k, ok := t.Udts[_type]; ok {
		return k.UDT
	}

	errorutils.Abort(errorutils.TypeError, errorutils.InvalidLLVMType, _type)
	return nil
}

// ImplicitTypeCast attempts to cast the given LLVM IR value to a specified target type.
//
// Parameters:
//
//	block  — the LLVM IR basic block where casting instructions are inserted.
//	target — the target type as a string (e.g., "i32", "float64", "string").
//	v      — the LLVM IR value to be cast.
//
// Returns:
//
//	value.Value — the resulting LLVM IR value after casting.
//	*ir.Block   — the (possibly updated) IR block reflecting the new instructions.
//
// Supported casts:
//   - Booleans: "boolean", "bool", "i1"
//   - Integers: "int8"/"i8", "int16"/"i16", "int32"/"i32", "int"/"int64"/"i64"
//   - Floats:   "float16"/"half", "float32"/"float", "float64"/"double"
//   - String:   "string", "i8*" (only from pointer types)
//   - Void:     "void" (produces nil)
//
// Additionally, user-defined types (UDTs) are resolved from the type registry.
// If the cast is invalid or unsupported, the function panics or aborts with a type error.
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
	case "array":
		if v.Type().Equal(ARRAYSTRUCT) {
			errorutils.Abort(errorutils.ImplicitTypeCastError, v.Type().String(), target)
		}
		return v, block
	}

	if k, ok := t.Udts[target]; ok {
		return ensureType(block, v, k.UDT), block
	}
	errorutils.Abort(errorutils.TypeError, errorutils.InvalidTargetType, target)
	return nil, block
}

// catchIntToIntDownCast inserts runtime checks for narrowing integer casts
// (downcasts) to detect overflow and raise an error if the value cannot fit
// in the destination integer type.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the source integer value to be downcast.
//	dst   — the target integer type (must be narrower than v’s type).
//
// Returns:
//
//	value.Value — the safely downcasted integer value, or a boolean when casting to i1.
//	*ir.Block   — the block after branching, pointing to the "safe" continuation path.
//
// Behavior:
//   - On overflow, a runtime error is raised, and execution is terminated via `unreachable`.
//   - On success, the value is truncated (`trunc`) to the destination type.
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

// catchFloatToIntDownCast inserts runtime checks for narrowing casts from
// floating-point values to integers, ensuring the value lies within the
// destination integer's bounds. If an overflow is detected, a runtime
// error is raised.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the floating-point value to be downcast.
//	dst   — the target integer type.
//
// Returns:
//
//	value.Value — the safely cast integer value (FP → SI).
//	*ir.Block   — the block after branching, pointing to the "safe" continuation path.
//
// Behavior:
//   - On overflow, a runtime error is raised and execution is terminated
//     with `unreachable`.
//   - On success, the float is converted to the destination integer type
//     using FPToSI.
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

// ImplicitIntCast casts a value to a target integer type, performing
// necessary runtime checks for overflows or width adjustments.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the value to be cast (integer or floating-point).
//	dst   — the destination integer type.
//
// Returns:
//
//	value.Value — the resulting LLVM IR value after casting.
//	*ir.Block   — the (possibly updated) block reflecting inserted instructions.
//
// Behavior:
//   - Boolean widening: i1 → larger integer uses ZExt; i1 → i1 returns unchanged.
//   - Integer upcast: smaller → larger integer uses SExt.
//   - Integer downcast: larger → smaller integer uses catchIntToIntDownCast
//     with overflow checks.
//   - Float → int: uses catchFloatToIntDownCast with overflow checks.
//   - Float → boolean: compares against 0.0 (non-zero → true).
//   - If the input type cannot be cast to an integer, the function aborts
//     with an implicit type cast error.
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
		errorutils.Abort(errorutils.ImplicitTypeCastError, v.Type().String(), "int")
	}

	if src.BitSize == 1 {
		if dst.BitSize == 1 {
			return v, b
		}
		// Boolean widen with ZExt
		return b.NewZExt(v, dst), b
	}

	// v is integer
	if src.BitSize > dst.BitSize {
		// downcast integer to integer (with overflow checks)
		return t.catchIntToIntDownCast(block, v, dst)
	}
	if src.BitSize < dst.BitSize {
		// extend smaller integer to larger
		return b.NewSExt(v, dst), b
	}
	return v, b
}

// catchFloatToFloatDowncast inserts runtime checks for floating-point
// narrowing casts (downcasts) to ensure that the value fits within the
// destination type's representable range. If an overflow occurs, a runtime
// error is raised.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the source floating-point value to be cast.
//	src   — the source floating-point type.
//	dst   — the destination floating-point type.
//
// Returns:
//
//	value.Value — the safely cast floating-point value.
//	*ir.Block   — the block after branching, pointing to the "safe" continuation path.
//
// Behavior:
//   - On overflow, a runtime error is raised and execution is terminated
//     with `unreachable`.
//   - In the safe path, the value is truncated (FPTrunc) to the destination type.
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

// catchIntToFloatDowncast inserts runtime checks for casting integers to
// floating-point types, ensuring the integer value fits within the
// representable range of the destination float. If the value exceeds the
// bounds, a runtime error is raised.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the integer value to be cast.
//	dst   — the destination floating-point type.
//
// Returns:
//
//	value.Value — the safely cast floating-point value.
//	*ir.Block   — the block after branching, pointing to the "safe" continuation path.
//
// Behavior:
//   - On overflow, a runtime error is raised and execution is terminated
//     with `unreachable`.
//   - In the safe path, the integer is converted to the requested float
//     type (SIToFP).
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

// ImplicitFloatCast casts a value to a target floating-point type, performing
// necessary runtime checks for safe conversions and width adjustments.
//
// Parameters:
//
//	block — the LLVM IR basic block where instructions are inserted.
//	v     — the value to be cast (floating-point or integer).
//	dst   — the destination floating-point type.
//
// Returns:
//
//	value.Value — the resulting LLVM IR value after casting.
//	*ir.Block   — the (possibly updated) block reflecting inserted instructions.
//
// Behavior:
//   - Float → float: uses catchFloatToFloatDowncast to handle upcasts and
//     downcasts with overflow checks.
//   - Integer → float: safe conversion with overflow checks.
//   - For i1, zero/one is promoted and converted to float.
//   - For larger integers, uses catchIntToFloatDowncast with overflow checks.
//   - If the input type cannot be cast to a float, the function panics.
func (t *TypeHandler) ImplicitFloatCast(block *ir.Block, v value.Value, dst *types.FloatType) (value.Value, *ir.Block) {
	switch src := v.Type().(type) {
	case *types.FloatType:
		return t.catchFloatToFloatDowncast(block, v, src, dst)

	case *types.IntType:
		// int -> float: special-case i1 -> treat as 0/1
		if src.BitSize == 1 {
			intVal := block.NewZExt(v, types.I8)
			floatVal := block.NewSIToFP(intVal, dst)
			return floatVal, block
		}
		return t.catchIntToFloatDowncast(block, v, dst)

	default:
		errorutils.Abort(errorutils.ImplicitTypeCastError, v.Type().String(), "float")
	}
	return nil, nil
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
