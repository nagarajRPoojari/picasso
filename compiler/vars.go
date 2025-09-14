package compiler

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// Var is a mutable variable (backed by an alloca slot).
type Var interface {
	Update(block *ir.Block, v value.Value)
	Load(block *ir.Block) value.Value
	Constant() constant.Constant
	Slot() value.Value
	Cast(block *ir.Block, v value.Value) value.Value
	Type() types.Type
}

type Boolean struct {
	NativeType *types.IntType
	Value      value.Value // pointer to i1
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

func (b *Boolean) Slot() value.Value { return b.Value }

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
func (c *Boolean) Type() types.Type { return c.NativeType }

/*** INT8 ***/
type Int8 struct {
	NativeType *types.IntType
	Value      value.Value
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
func (i *Int8) Slot() value.Value                     { return i.Value }
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
func (c *Int8) Type() types.Type { return c.NativeType }

/*** INT16 ***/
type Int16 struct {
	NativeType *types.IntType
	Value      value.Value
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
func (i *Int16) Slot() value.Value                     { return i.Value }
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
func (c *Int16) Type() types.Type { return c.NativeType }

/*** INT32 ***/
type Int32 struct {
	NativeType *types.IntType
	Value      value.Value
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
func (i *Int32) Slot() value.Value                     { return i.Value }
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
func (c *Int32) Type() types.Type { return c.NativeType }

/*** INT64 ***/
type Int64 struct {
	NativeType *types.IntType
	Value      value.Value
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
func (i *Int64) Slot() value.Value                     { return i.Value }
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
func (c *Int64) Type() types.Type { return c.NativeType }

/*** FLOAT16 (half) ***/
type Float16 struct {
	NativeType *types.FloatType
	Value      value.Value
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
func (f *Float16) Slot() value.Value { return f.Value }
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
func (c *Float16) Type() types.Type { return c.NativeType }

type Float32 struct {
	NativeType *types.FloatType
	Value      value.Value
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
func (f *Float32) Slot() value.Value { return f.Value }
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
func (c *Float32) Type() types.Type { return c.NativeType }

/*** FLOAT64 (double) ***/
type Float64 struct {
	NativeType *types.FloatType
	Value      value.Value
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
func (f *Float64) Slot() value.Value                     { return f.Value }
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
func (c *Float64) Type() types.Type { return c.NativeType }

type Class struct {
	Name  string
	UDT   types.Type
	Value value.Value
}

func NewClass(block *ir.Block, name string, udt types.Type) *Class {
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
func (s *Class) UpdateField(block *ir.Block, idx int, v value.Value, expected types.Type) {
	val := ensureType(block, v, expected)
	fieldPtr := s.FieldPtr(block, idx)
	block.NewStore(val, fieldPtr)
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
func (s *Class) Slot() value.Value { return s.Value }

func (c *Class) Type() types.Type { return c.UDT }

func ensureType(block *ir.Block, v value.Value, target types.Type) value.Value {
	if v.Type().Equal(target) {
		return v
	}

	switch src := v.Type().(type) {
	case *types.IntType:
		dst, ok := target.(*types.IntType)
		if !ok {
			break
		}
		if src.BitSize < dst.BitSize {
			return block.NewZExt(v, dst) // or NewSExt if signed
		} else if src.BitSize > dst.BitSize {
			return block.NewTrunc(v, dst)
		}
		return v

	case *types.FloatType:
		dst, ok := target.(*types.FloatType)
		if !ok {
			break
		}
		if src.Kind == dst.Kind {
			return v
		}
		// Promote/demote based on known float kinds
		if floatRank(src.Kind) < floatRank(dst.Kind) {
			return block.NewFPExt(v, dst) // promote
		}
	}

	// fallback: bitcast (pointers etc.)
	return block.NewBitCast(v, target)
}

type Null struct {
	Target types.Type
	Value  *ir.InstAlloca
}

func NewNull() *Null {
	return &Null{}
}

func (n *Null) Load(block *ir.Block) value.Value {
	return block.NewLoad(n.Target, n.Value)
}

func (n *Null) Update(block *ir.Block, v value.Value) {
	if !v.Type().Equal(n.Target) {
		panic(fmt.Sprintf("cannot assign %s to null of type %s", v.Type(), n.Target))
	}
	block.NewStore(v, n.Value)
}

func (n *Null) Type() types.Type {
	return n.Target
}

func (n *Null) Cast(block *ir.Block, v value.Value) value.Value {
	panic(fmt.Sprintf("cannot cast %s to null type %s", v.Type(), n.Target))
}

func (n *Null) GoValue() interface{} {
	return nil
}
