package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/compiler/gc"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

// Class is a custom user defined data type
type Class struct {
	Name string      // class name (for lookup)
	UDT  types.Type  // always pointer-to-struct type (types.PointerType)
	Ptr  value.Value // pointer value (pointer-to-struct, i.e. the object address)
}

func NewClass(block *ir.Block, name string, udt types.Type) *Class {
	// Normalize udt so s.UDT is always *types.PointerType (pointer-to-struct)
	var ptrType *types.PointerType
	switch t := udt.(type) {
	case *types.StructType:
		ptrType = types.NewPointer(t)
	case *types.PointerType:
		// if pointer already, ensure it points to a struct (optional check)
		if _, ok := t.ElemType.(*types.StructType); !ok {
			panic(fmt.Sprintf("NewClass expects pointer-to-struct or struct, got pointer to %T", t.ElemType))
		}
		ptrType = t
	default:
		panic(fmt.Sprintf("NewClass expects struct or pointer-to-struct, got %T", udt))
	}

	// compute allocation size: gep on element type (struct) with null pointer and index 1
	zero := constant.NewNull(ptrType)
	one := constant.NewInt(types.I32, 1)
	// Get element ptr MUST use the element (struct) type
	gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
	size := constant.NewPtrToInt(gep, types.I64)

	// Call GC allocator
	mem := block.NewCall(gc.Instance.Alloc(), size)

	// Bitcast to your struct pointer type
	ptr := block.NewBitCast(mem, ptrType)

	return &Class{
		Name: name,
		UDT:  ptrType, // store pointer-to-struct type
		Ptr:  ptr,     // pointer to allocated object
	}
}

func (s *Class) Update(block *ir.Block, v value.Value) {
	if v == nil {
		errorsx.PanicCompilationError(fmt.Sprintf("cannot update object with nil value: %v", v))
	}

	// Ensure s.UDT is pointer type
	sPtr, ok := s.UDT.(*types.PointerType)
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("internal error: Class.UDT is not a pointer type: %T", s.UDT))
	}

	// Case: v is pointer-to-struct and matches s.UDT
	if pv, ok := v.Type().(*types.PointerType); ok && pv.Equal(sPtr) {
		// load the struct value from v and store the struct into the object's address (s.Ptr)
		s.Ptr = v
		return
	}

	// Case: v is the struct value itself (elem type)
	if v.Type().Equal(sPtr.ElemType) {
		// store the struct value into the object's address
		block.NewStore(v, s.Ptr)
		return
	}

	// Fallback: try to convert/cast to pointer-to-struct or struct appropriately, then store
	// If we can obtain a struct value, store it; otherwise try to bitcast ptr and load.
	val := ensureType(block, v, sPtr.ElemType) // try to get struct value
	block.NewStore(val, s.Ptr)
}

func (s *Class) Load(block *ir.Block) value.Value {
	// Return the struct value loaded from object's address (s.Ptr).
	// s.UDT is pointer-to-struct, so we must load ElemType.
	if _, ok := s.UDT.(*types.PointerType); ok {
		return s.Ptr
	}
	errorsx.PanicCompilationError(fmt.Sprintf("internal error: Class.UDT is not pointer type: %T", s.UDT))
	return nil
}

func (s *Class) FieldPtr(block *ir.Block, idx int) value.Value {
	zero := constant.NewInt(types.I32, 0)
	i := constant.NewInt(types.I32, int64(idx))

	// unwrap pointer-to-struct
	sPtr, ok := s.UDT.(*types.PointerType)
	if !ok {
		errorsx.PanicCompilationError(fmt.Sprintf("internal error: Class.UDT is not pointer type: %T", s.UDT))
	}
	elem := sPtr.ElemType // struct type

	// GEP: base is the pointer-to-struct (object address) and we index into the struct
	return block.NewGetElementPtr(elem, s.Ptr, zero, i)
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

	sPtr, ok := s.UDT.(*types.PointerType)
	if !ok {
		return nil, fmt.Errorf("internal error: Class.UDT is not pointer type: %T", s.UDT)
	}

	// If already pointer to struct (the desired pointer type)
	if v.Type().Equal(sPtr) {
		return v, nil
	}

	// If we have struct value, alloca+store and return pointer to it
	if v.Type().Equal(sPtr.ElemType) {
		tmp := block.NewAlloca(v.Type())
		block.NewStore(v, tmp)
		return tmp, nil
	}

	// If any pointer, bitcast to this pointer type
	if _, ok := v.Type().(*types.PointerType); ok {
		return block.NewBitCast(v, sPtr), nil
	}

	// Fallback: try to convert value to struct type then alloca
	val := ensureType(block, v, sPtr.ElemType)
	tmp := block.NewAlloca(val.Type())
	block.NewStore(val, tmp)
	return tmp, nil
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
		switch dst := target.(type) {
		case *types.IntType:
			if src.BitSize < dst.BitSize {
				return block.NewZExt(v, dst)
			} else if src.BitSize > dst.BitSize {
				return block.NewTrunc(v, dst)
			}
			return v
		case *types.FloatType:
			return block.NewSIToFP(v, dst)
		}
	case *types.FloatType:
		switch dst := target.(type) {
		case *types.FloatType:
			if floatRank(src.Kind) < floatRank(dst.Kind) {
				return block.NewFPExt(v, dst)
			} else if floatRank(src.Kind) > floatRank(dst.Kind) {
				return block.NewFPTrunc(v, dst)
			}
			return v
		case *types.IntType:
			return block.NewFPToSI(v, dst)
		}
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		if _, ok2 := target.(*types.PointerType); ok2 {
			return block.NewBitCast(v, target)
		}
	}

	return block.NewBitCast(v, target)
}

func (f *Class) NativeTypeString() string { return f.Name }
