package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Var is a mutable variable (backed by an alloca slot).
type Var interface {
	Update(block *ir.Block, v value.Value)
	Load(block *ir.Block) value.Value
	Constant() constant.Constant
	Slot() *ir.InstAlloca
	Cast(block *ir.Block, v value.Value) value.Value
}

/*** BOOLEAN ***/
type Boolean struct {
	NativeType *types.IntType
	Value      *ir.InstAlloca // pointer to i1
	GoVal      bool
}

func NewBooleanVar(block *ir.Block, init bool) *Boolean {
	slot := block.NewAlloca(types.I1)
	block.NewStore(constant.NewInt(types.I1, btoi(init)), slot)
	return &Boolean{NativeType: types.I1, Value: slot, GoVal: init}
}

func (b *Boolean) Update(block *ir.Block, v value.Value) {
	block.NewStore(v, b.Value)
}

func (b *Boolean) Load(block *ir.Block) value.Value {
	return block.NewLoad(types.I1, b.Value)
}

func (b *Boolean) Constant() constant.Constant {
	return constant.NewInt(types.I1, btoi(b.GoVal))
}

func (b *Boolean) Slot() *ir.InstAlloca { return b.Value }

func (b *Boolean) Cast(block *ir.Block, v value.Value) value.Value {
	switch v.Type().(type) {
	case *types.IntType:
		// cast any int to i1
		if v.Type().(*types.IntType).BitSize == 1 {
			return v
		}
		return block.NewTrunc(v, types.I1)
	default:
		panic(fmt.Sprintf("cannot cast %s to boolean", v.Type()))
	}
}

/*** INT8 ***/
type Int8 struct {
	NativeType *types.IntType
	Value      *ir.InstAlloca
	GoVal      int8
}

func NewInt8Var(block *ir.Block, init int8) *Int8 {
	slot := block.NewAlloca(types.I8)
	block.NewStore(constant.NewInt(types.I8, int64(init)), slot)
	return &Int8{NativeType: types.I8, Value: slot, GoVal: init}
}

func (i *Int8) Update(block *ir.Block, v value.Value) { block.NewStore(v, i.Value) }
func (i *Int8) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I8, i.Value) }
func (i *Int8) Constant() constant.Constant           { return constant.NewInt(types.I8, int64(i.GoVal)) }
func (i *Int8) Slot() *ir.InstAlloca                  { return i.Value }
func (i *Int8) Cast(block *ir.Block, v value.Value) value.Value {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 8 {
			return block.NewTrunc(v, types.I8)
		} else if t.BitSize < 8 {
			return block.NewSExt(v, types.I8)
		}
		return v
	case *types.FloatType:
		return block.NewFPToSI(v, types.I8)
	default:
		panic(fmt.Sprintf("cannot cast %s to int8", v.Type()))
	}
}

/*** INT16 ***/
type Int16 struct {
	NativeType *types.IntType
	Value      *ir.InstAlloca
	GoVal      int16
}

func NewInt16Var(block *ir.Block, init int16) *Int16 {
	slot := block.NewAlloca(types.I16)
	block.NewStore(constant.NewInt(types.I16, int64(init)), slot)
	return &Int16{NativeType: types.I16, Value: slot, GoVal: init}
}

func (i *Int16) Update(block *ir.Block, v value.Value) { block.NewStore(v, i.Value) }
func (i *Int16) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I16, i.Value) }
func (i *Int16) Constant() constant.Constant           { return constant.NewInt(types.I16, int64(i.GoVal)) }
func (i *Int16) Slot() *ir.InstAlloca                  { return i.Value }
func (i *Int16) Cast(block *ir.Block, v value.Value) value.Value {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 16 {
			return block.NewTrunc(v, types.I16)
		} else if t.BitSize < 16 {
			return block.NewSExt(v, types.I16)
		}
		return v
	case *types.FloatType:
		return block.NewFPToSI(v, types.I16)
	default:
		panic(fmt.Sprintf("cannot cast %s to int16", v.Type()))
	}
}

/*** INT32 ***/
type Int32 struct {
	NativeType *types.IntType
	Value      *ir.InstAlloca
	GoVal      int32
}

func NewInt32Var(block *ir.Block, init int32) *Int32 {
	slot := block.NewAlloca(types.I32)
	block.NewStore(constant.NewInt(types.I32, int64(init)), slot)
	return &Int32{NativeType: types.I32, Value: slot, GoVal: init}
}

func (i *Int32) Update(block *ir.Block, v value.Value) { block.NewStore(v, i.Value) }
func (i *Int32) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I32, i.Value) }
func (i *Int32) Constant() constant.Constant           { return constant.NewInt(types.I32, int64(i.GoVal)) }
func (i *Int32) Slot() *ir.InstAlloca                  { return i.Value }
func (i *Int32) Cast(block *ir.Block, v value.Value) value.Value {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 32 {
			return block.NewTrunc(v, types.I32)
		} else if t.BitSize < 32 {
			return block.NewSExt(v, types.I32)
		}
		return v
	case *types.FloatType:
		return block.NewFPToSI(v, types.I32)
	default:
		panic(fmt.Sprintf("cannot cast %s to int32", v.Type()))
	}
}

/*** INT64 ***/
type Int64 struct {
	NativeType *types.IntType
	Value      *ir.InstAlloca
	GoVal      int64
}

func NewInt64Var(block *ir.Block, init int64) *Int64 {
	slot := block.NewAlloca(types.I64)
	block.NewStore(constant.NewInt(types.I64, init), slot)
	return &Int64{NativeType: types.I64, Value: slot, GoVal: init}
}

func (i *Int64) Update(block *ir.Block, v value.Value) { block.NewStore(v, i.Value) }
func (i *Int64) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I64, i.Value) }
func (i *Int64) Constant() constant.Constant           { return constant.NewInt(types.I64, i.GoVal) }
func (i *Int64) Slot() *ir.InstAlloca                  { return i.Value }
func (i *Int64) Cast(block *ir.Block, v value.Value) value.Value {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 64 {
			return block.NewTrunc(v, types.I64)
		} else if t.BitSize < 64 {
			return block.NewSExt(v, types.I64)
		}
		return v
	case *types.FloatType:
		return block.NewFPToSI(v, types.I64)
	default:
		panic(fmt.Sprintf("cannot cast %s to int64", v.Type()))
	}
}

/*** FLOAT16 (half) ***/
type Float16 struct {
	NativeType *types.FloatType
	Value      *ir.InstAlloca
	GoVal      float32
}

func NewFloat16Var(block *ir.Block, init float32) *Float16 {
	slot := block.NewAlloca(types.Half)
	block.NewStore(constant.NewFloat(types.Half, float64(init)), slot)
	return &Float16{NativeType: types.Half, Value: slot, GoVal: init}
}

func (f *Float16) Update(block *ir.Block, v value.Value) { block.NewStore(v, f.Value) }
func (f *Float16) Load(block *ir.Block) value.Value      { return block.NewLoad(types.Half, f.Value) }
func (f *Float16) Constant() constant.Constant {
	return constant.NewFloat(types.Half, float64(f.GoVal))
}
func (f *Float16) Slot() *ir.InstAlloca { return f.Value }
func (f *Float16) Cast(block *ir.Block, v value.Value) value.Value {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Half)
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.NewFPTrunc(v, types.Half)
		case types.Half:
			return block.NewFPExt(v, types.Half)
		case types.Float:
			return v
		}
	}
	panic(fmt.Sprintf("cannot cast %s to float16", v.Type()))
}

/*** FLOAT32 ***/
type Float32 struct {
	NativeType *types.FloatType
	Value      *ir.InstAlloca
	GoVal      float32
}

func NewFloat32Var(block *ir.Block, init float32) *Float32 {
	slot := block.NewAlloca(types.Float)
	block.NewStore(constant.NewFloat(types.Float, float64(init)), slot)
	return &Float32{NativeType: types.Float, Value: slot, GoVal: init}
}

func (f *Float32) Update(block *ir.Block, v value.Value) { block.NewStore(v, f.Value) }
func (f *Float32) Load(block *ir.Block) value.Value      { return block.NewLoad(types.Float, f.Value) }
func (f *Float32) Constant() constant.Constant {
	return constant.NewFloat(types.Float, float64(f.GoVal))
}
func (f *Float32) Slot() *ir.InstAlloca { return f.Value }
func (f *Float32) Cast(block *ir.Block, v value.Value) value.Value {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Float)
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.NewFPTrunc(v, types.Float)
		case types.Half:
			return block.NewFPExt(v, types.Float)
		case types.Float:
			return v
		}
	}
	panic(fmt.Sprintf("cannot cast %s to float32", v.Type()))
}

/*** FLOAT64 (double) ***/
type Float64 struct {
	NativeType *types.FloatType
	Value      *ir.InstAlloca
	GoVal      float64
}

func NewFloat64Var(block *ir.Block, init float64) *Float64 {
	slot := block.NewAlloca(types.Double)
	block.NewStore(constant.NewFloat(types.Double, init), slot)
	return &Float64{NativeType: types.Double, Value: slot, GoVal: init}
}

func (f *Float64) Update(block *ir.Block, v value.Value) { block.NewStore(v, f.Value) }
func (f *Float64) Load(block *ir.Block) value.Value      { return block.NewLoad(types.Double, f.Value) }
func (f *Float64) Constant() constant.Constant           { return constant.NewFloat(types.Double, f.GoVal) }
func (f *Float64) Slot() *ir.InstAlloca                  { return f.Value }
func (f *Float64) Cast(block *ir.Block, v value.Value) value.Value {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Double)
	case *types.FloatType:
		switch v.Type() {
		case types.Float:
			return block.NewFPExt(v, types.Double)
		case types.Half:
			return block.NewFPExt(v, types.Double)
		case types.Double:
			return v
		}
	}
	panic(fmt.Sprintf("cannot cast %s to float64", v.Type()))
}

/*** STRING (as i8*) ***/
// Note: this treats the variable as an i8* pointer. You may want a more
// elaborate representation (global char array + GEP) for literals.
type String struct {
	NativeType *types.PointerType
	Value      *ir.InstAlloca // pointer to i8*
	GoVal      string
}

func NewStringVar(block *ir.Block, init string) *String {
	ptrTy := types.NewPointer(types.I8)
	slot := block.NewAlloca(ptrTy)
	// For simplicity we store a null pointer if no initializer is desired.
	// To initialize with a literal you'd normally create a global char array and store its pointer.
	if init != "" {
		// Not initializing to the literal pointer here; caller should set after creating a global literal.
		block.NewStore(constant.NewNull(ptrTy), slot)
	} else {
		block.NewStore(constant.NewNull(ptrTy), slot)
	}
	return &String{NativeType: ptrTy, Value: slot, GoVal: init}
}

func (s *String) Update(block *ir.Block, v value.Value) { block.NewStore(v, s.Value) }
func (s *String) Load(block *ir.Block) value.Value      { return block.NewLoad(s.NativeType, s.Value) }
func (s *String) Constant() constant.Constant {
	return constant.NewCharArrayFromString(s.GoVal)
}
func (s *String) Slot() *ir.InstAlloca { return s.Value }

type Class struct {
	Name  string
	UDT   types.Type
	Value *ir.InstAlloca
}

func NewClass(block *ir.Block, name string, udt types.Type) *Class {
	// Define struct type
	ptr := block.NewAlloca(udt)

	return &Class{
		Name:  name,
		UDT:   udt,
		Value: ptr,
	}
}

// Store a whole struct value into the alloca
func (s *Class) Update(block *ir.Block, v value.Value) {
	block.NewStore(v, s.Value)
}

// Load entire struct value from alloca
func (s *Class) Load(block *ir.Block) value.Value {
	return block.NewLoad(s.UDT, s.Value)
}

// Access a field by index
func (s *Class) FieldPtr(block *ir.Block, idx int) value.Value {
	return block.NewGetElementPtr(s.UDT, s.Value, constant.NewInt(types.I32, 0), constant.NewInt(types.I32, int64(idx)))
}

// Update a single field
func (s *Class) UpdateField(block *ir.Block, idx int, v value.Value) {
	fieldPtr := s.FieldPtr(block, idx)
	block.NewStore(v, fieldPtr)
}

// Load a single field
func (s *Class) LoadField(block *ir.Block, idx int, fieldType types.Type) value.Value {
	fieldPtr := s.FieldPtr(block, idx)
	return block.NewLoad(fieldType, fieldPtr)
}

func (f *Class) Cast(block *ir.Block, v value.Value) value.Value {
	return nil
}

func (f *Class) Constant() constant.Constant {
	return nil
}
func (s *Class) Slot() *ir.InstAlloca { return s.Value }

// CastToType takes a target type name (e.g. "float64", "int8")
// and a value, and emits the appropriate cast instruction in `block`.
func CastToType(block *ir.Block, target string, v value.Value) value.Value {
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
		return intCast(block, v, types.I8)
	case "int16", "i16":
		return intCast(block, v, types.I16)
	case "int32", "i32":
		return intCast(block, v, types.I32)
	case "int", "int64", "i64":
		return intCast(block, v, types.I64)

	case "float16", "half":
		return floatCast(block, v, types.Half)
	case "float32", "float":
		return floatCast(block, v, types.Float)
	case "float64", "double":
		return floatCast(block, v, types.Double)

	default:
		panic("unsupported target type: " + target)
	}
}

func intCast(block *ir.Block, v value.Value, dst *types.IntType) value.Value {
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

func floatCast(block *ir.Block, v value.Value, dst *types.FloatType) value.Value {
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
