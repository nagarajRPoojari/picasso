package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	errorutils "github.com/nagarajRPoojari/x-lang/generator/error"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

// Class is a custom user defined data type
type Class struct {
	Name string      // class name (for lookup)
	UDT  types.Type  // always pointer-to-struct type (types.PointerType)
	Ptr  value.Value // pointer value (pointer-to-struct, i.e. the object address)
}

func NewClass(block *bc.BlockHolder, name string, udt types.Type) *Class {
	// Normalize UDT so it's always pointer-to-struct (*types.PointerType)
	var ptrType *types.PointerType
	switch t := udt.(type) {
	case *types.StructType:
		ptrType = types.NewPointer(t)
	case *types.PointerType:
		if _, ok := t.ElemType.(*types.StructType); !ok {
			panic(fmt.Sprintf("NewClass expects pointer-to-struct or struct, got pointer to %T", t.ElemType))
		}
		ptrType = t
	default:
		panic(fmt.Sprintf("NewClass expects struct or pointer-to-struct, got %T", udt))
	}

	// === Allocate heap memory for struct (runtime allocation) ===
	zero := constant.NewNull(ptrType)
	one := constant.NewInt(types.I32, 1)
	gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
	size := constant.NewPtrToInt(gep, types.I64)

	// Call allocator (runtime malloc/GC alloc)
	mem := block.N.NewCall(c.Instance.Funcs[c.ALLOC], size)

	// Bitcast to your struct pointer type
	ptr := block.N.NewBitCast(mem, ptrType)

	// === Create stack slot to store the pointer ===
	// alloca type = pointer to (pointer-to-struct)
	alloca := block.N.NewAlloca(ptrType)
	block.N.NewStore(ptr, alloca)

	// s.Ptr now points to a memory cell that *contains* the pointer
	return &Class{
		Name: name,
		UDT:  ptrType, // pointer-to-struct type
		Ptr:  alloca,  // alloca holds the runtime pointer value
	}
}

func (s *Class) Update(bh *bc.BlockHolder, v value.Value) {
	block := bh.N

	if v == nil {
		errorutils.Abort(errorutils.InternalError, "cannot update object with nil value")
	}

	ptrType, ok := s.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, fmt.Sprintf("Class.UDT is not pointer-to-struct: %T", s.UDT))
	}

	// Ensure s.Ptr is allocated (holds the pointer to the struct)
	if s.Ptr == nil {
		// Allocate memory for the pointer itself (stack slot)
		s.Ptr = block.NewAlloca(ptrType)
		// Also allocate memory for the struct instance
		// zero := constant.NewNull(ptrType)
		// one := constant.NewInt(types.I32, 1)
		// gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
		// size := constant.NewPtrToInt(gep, types.I64)
		// mem := block.NewCall(c.Instance.Funcs[c.ALLOC], size)
		// structPtr := block.NewBitCast(mem, ptrType)
		// block.NewStore(structPtr, s.Ptr)
	}

	// Case 1: v already matches our pointer-to-struct type
	if pv, ok := v.Type().(*types.PointerType); ok && pv.Equal(ptrType) {
		block.NewStore(v, s.Ptr)
		return
	}

	// Case 2: v is a raw struct, not pointer – convert it
	val := ensureType(bh, v, ptrType.ElemType)
	// Load current struct pointer
	currPtr := block.NewLoad(ptrType, s.Ptr)
	block.NewStore(val, currPtr)
}

func (a *Class) UpdateV2(block *bc.BlockHolder, v *Class) {
	*a = *v
}

func (s *Class) Load(bh *bc.BlockHolder) value.Value {
	block := bh.N

	sPtr, ok := s.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, fmt.Sprintf("Class.UDT is not pointer type: %T", s.UDT))
		return nil
	}

	// load runtime pointer-to-struct from stack slot
	ptrVal := block.NewLoad(sPtr, s.Ptr)
	return ptrVal
}

func (s *Class) FieldPtr(block *bc.BlockHolder, idx int) value.Value {
	zero := constant.NewInt(types.I32, 0)
	i := constant.NewInt(types.I32, int64(idx))

	// unwrap pointer-to-struct
	sPtr, ok := s.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, fmt.Sprintf("Class.UDT is not a pointer type: %T", s.UDT))
	}
	elem := sPtr.ElemType // struct type

	// GEP: base is the pointer-to-struct (object address) and we index into the struct

	return block.N.NewGetElementPtr(elem, s.Load(block), zero, i)
}

func (s *Class) UpdateField(bh *bc.BlockHolder, idx int, v value.Value, expected types.Type) {
	block := bh.N
	if v == nil {
		// errorsx.PanicCompilationError(fmt.Sprintf("cannot update field with nil value: %v", v))
		return
	}
	val := ensureType(bh, v, expected)
	fieldPtr := s.FieldPtr(bh, idx)
	block.NewStore(val, fieldPtr)
}

func (s *Class) LoadField(block *bc.BlockHolder, idx int, fieldType types.Type) value.Value {
	fieldPtr := s.FieldPtr(block, idx)
	return block.N.NewLoad(fieldType, fieldPtr)
}

func (s *Class) Cast(bh *bc.BlockHolder, v value.Value) (value.Value, error) {
	block := bh.N
	if v == nil {
		errorutils.Abort(errorutils.InternalError, fmt.Sprintf("cannot cast nil value: %v", v))
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
	val := ensureType(bh, v, sPtr.ElemType)
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

func ensureType(block *bc.BlockHolder, v value.Value, target types.Type) value.Value {
	srcType := v.Type()

	if srcType.Equal(target) {
		return v
	}

	// Handle pointer types
	srcPtr, okSrc := srcType.(*types.PointerType)
	tgtPtr, okTgt := target.(*types.PointerType)

	if okSrc && okTgt {
		srcStruct, srcIsStruct := srcPtr.ElemType.(*types.StructType)
		tgtStruct, tgtIsStruct := tgtPtr.ElemType.(*types.StructType)

		if srcIsStruct && tgtIsStruct {
			// Allow empty struct to be cast to any struct
			if len(srcStruct.Fields) == 0 {
				return block.N.NewBitCast(v, target)
			}

			// Child -> Parent layout check
			if len(srcStruct.Fields) < len(tgtStruct.Fields) {
				panic("ensureType: source struct has fewer fields than target parent")
			}

			for i := range tgtStruct.Fields {
				if !srcStruct.Fields[i].Equal(tgtStruct.Fields[i]) {
					panic(fmt.Sprintf("ensureType: field %d type mismatch: %v vs %v", i, srcStruct.Fields[i], tgtStruct.Fields[i]))
				}
			}

			// Safe to bitcast child pointer to parent pointer
			return block.N.NewBitCast(v, target)
		}

		// Not structs, fallback: generic pointer bitcast
		return block.N.NewBitCast(v, target)
	}

	// Bitcast function pointers
	if _, okSrcFunc := srcType.(*types.FuncType); okSrcFunc {
		if _, okTgtFunc := target.(*types.FuncType); okTgtFunc {
			return block.N.NewBitCast(v, target)
		}
	}
	panic(fmt.Sprintf("ensureType: cannot convert %v -> %v", srcType, target))
}

func (f *Class) NativeTypeString() string { return f.Name }
