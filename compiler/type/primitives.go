package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/compiler/utils"

	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Boolean hold single bit of information
type Boolean struct {
	NativeType *types.IntType
	Value      value.Value // pointer to i1
	GoVal      bool
}

func NewBooleanVar(block *ir.Block, init bool) *Boolean {
	slot := block.NewAlloca(types.I1)
	block.NewStore(constant.NewInt(types.I1, utils.BtoI(init)), slot)
	return &Boolean{NativeType: types.I1, Value: slot, GoVal: init}
}
func (b *Boolean) Update(block *ir.Block, v value.Value) { block.NewStore(v, b.Value) }
func (b *Boolean) Load(block *ir.Block) value.Value      { return block.NewLoad(types.I1, b.Value) }
func (b *Boolean) Constant() constant.Constant           { return constant.NewInt(types.I1, utils.BtoI(b.GoVal)) }
func (b *Boolean) Slot() value.Value                     { return b.Value }
func (c *Boolean) Type() types.Type                      { return c.NativeType }
func (b *Boolean) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		// cast any int to i1
		if v.Type().(*types.IntType).BitSize == 1 {
			return v, nil
		}
		// if incoming int type is other than types.I1, truncate it to single bit
		return block.NewTrunc(v, types.I1), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to booolean", v))
	}
}

// Int8 hold 1 byte of information
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
func (c *Int8) Type() types.Type                      { return c.NativeType }
func (i *Int8) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		// truncate to 1 byte if byte size if greater
		if t.BitSize > 8 {
			return block.NewTrunc(v, types.I8), nil
		} else if t.BitSize < 8 {
			return block.NewSExt(v, types.I8), nil
		}
		return v, nil
	case *types.FloatType:
		return block.NewFPToSI(v, types.I8), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int7", v))
	}
}

// Int16 store 2 bytes of information
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
func (c *Int16) Type() types.Type                      { return c.NativeType }
func (i *Int16) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 16 {
			return block.NewTrunc(v, types.I16), nil
		} else if t.BitSize < 16 {
			return block.NewSExt(v, types.I16), nil
		}
		return v, nil
	case *types.FloatType:
		return block.NewFPToSI(v, types.I16), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int16", v))

	}
}

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
func (i *Int32) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 32 {
			return block.NewTrunc(v, types.I32), nil
		} else if t.BitSize < 32 {
			return block.NewSExt(v, types.I32), nil
		}
		return v, nil
	case *types.FloatType:
		return block.NewFPToSI(v, types.I32), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int32", v))
	}
}
func (c *Int32) Type() types.Type { return c.NativeType }

// Int64 stores 8 bytes int value
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
func (i *Int64) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch t := v.Type().(type) {
	case *types.IntType:
		if t.BitSize > 64 {
			return block.NewTrunc(v, types.I64), nil
		} else if t.BitSize < 64 {
			return block.NewSExt(v, types.I64), nil
		}
		return v, nil
	case *types.FloatType:
		return block.NewFPToSI(v, types.I64), nil
	default:
		return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
	}
}
func (c *Int64) Type() types.Type { return c.NativeType }

// Float16/Half stores 2 byte floating point value
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
func (f *Float16) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Half), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.NewFPTrunc(v, types.Half), nil
		case types.Half:
			return block.NewFPExt(v, types.Half), nil
		case types.Float:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to int64", v))
}
func (c *Float16) Type() types.Type { return c.NativeType }

// Float32 stores 4 byte floating point value
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
func (c *Float32) Type() types.Type  { return c.NativeType }
func (f *Float32) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Float), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Double:
			return block.NewFPTrunc(v, types.Float), nil
		case types.Half:
			return block.NewFPExt(v, types.Float), nil
		case types.Float:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to float32", v))
}

// Float64 stores 8 byte floating point value
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
func (c *Float64) Type() types.Type                      { return c.NativeType }
func (f *Float64) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	switch v.Type().(type) {
	case *types.IntType:
		return block.NewSIToFP(v, types.Double), nil
	case *types.FloatType:
		switch v.Type() {
		case types.Float:
			return block.NewFPExt(v, types.Double), nil
		case types.Half:
			return block.NewFPExt(v, types.Double), nil
		case types.Double:
			return v, nil
		}
	}
	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to float64", v))
}
