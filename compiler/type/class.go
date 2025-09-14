package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Class is a custom user defined data type
type Class struct {
	Name string      // class name (for lookup)
	UDT  types.Type  // the struct type
	Ptr  value.Value // pointer to the struct (alloca, GEP, etc.)
}

func NewClass(block *ir.Block, name string, udt types.Type) *Class {
	ptr := block.NewAlloca(udt)
	return &Class{
		Name: name,
		UDT:  udt,
		Ptr:  ptr,
	}
}

func (s *Class) Update(block *ir.Block, v value.Value) {
	if v == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("cannot update object with nil value: %v", v))
	}

	if ptrT, ok := v.Type().(*types.PointerType); ok {
		if ptrT.Equal(s.UDT) {
			val := block.NewLoad(s.UDT, v)
			block.NewStore(val, s.Ptr)
			return
		}
	}

	if v.Type().Equal(s.UDT) {
		block.NewStore(v, s.Ptr)
		return
	}

	val := ensureType(block, v, s.UDT)
	block.NewStore(val, s.Ptr)
}

func (s *Class) Load(block *ir.Block) value.Value {
	return block.NewLoad(s.UDT, s.Ptr)
}

func (s *Class) FieldPtr(block *ir.Block, idx int) value.Value {
	zero := constant.NewInt(types.I32, 0)
	i := constant.NewInt(types.I32, int64(idx))
	return block.NewGetElementPtr(s.UDT, s.Ptr, zero, i)
}

func (s *Class) UpdateField(block *ir.Block, idx int, v value.Value, expected types.Type) {
	if v == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("cannot update field with nil value: %v", v))
	}
	val := ensureType(block, v, expected)
	fieldPtr := s.FieldPtr(block, idx)
	block.NewStore(val, fieldPtr)
}

func (s *Class) LoadField(block *ir.Block, idx int, fieldType types.Type) value.Value {
	fieldPtr := s.FieldPtr(block, idx)
	return block.NewLoad(fieldType, fieldPtr)
}

func (s *Class) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	if v == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("cannot cast nil value: %v", v))
	}

	// If already the struct value
	if v.Type().Equal(s.UDT) {
		return v, nil
	}

	// If pointer to struct, load value
	if ptrT, ok := v.Type().(*types.PointerType); ok && ptrT.Equal(s.UDT) {
		return block.NewLoad(s.UDT, v), nil
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		ptrToUDT := types.NewPointer(s.UDT)
		return block.NewBitCast(v, ptrToUDT), nil
	}

	return ensureType(block, v, s.UDT), nil
}

func (s *Class) Constant() constant.Constant {
	return nil
}

func (s *Class) Slot() value.Value {
	return s.Ptr
}

func (s *Class) Type() types.Type {
	return s.UDT
}

func ensureType(block *ir.Block, v value.Value, target types.Type) value.Value {
	if v == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("cannot determine type of nil value: %v", v))
	}
	if v.Type().Equal(target) {
		return v
	}

	switch src := v.Type().(type) {
	case *types.IntType:
		if dst, ok := target.(*types.IntType); ok {
			// promote/demote by bit-size
			if src.BitSize < dst.BitSize {
				return block.NewZExt(v, dst)
			} else if src.BitSize > dst.BitSize {
				return block.NewTrunc(v, dst)
			}
			return v
		}
	case *types.FloatType:
		if dst, ok := target.(*types.FloatType); ok {
			if src.Kind == dst.Kind {
				return v
			}
			if floatRank(src.Kind) < floatRank(dst.Kind) {
				return block.NewFPExt(v, dst)
			}
		}
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		if _, ok2 := target.(*types.PointerType); ok2 {
			return block.NewBitCast(v, target)
		}
	}

	return block.NewBitCast(v, target)
}
