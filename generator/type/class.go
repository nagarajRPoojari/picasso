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

type Class struct {
	Name string             // Name for debugging/lookup
	UDT  *types.PointerType // The type of the object (e.g., %MyStruct*)
	Ptr  value.Value        // The stack slot (alloca) holding the pointer (e.g., %MyStruct**)
}

// NewClass creates a new object instance.
// 1. Allocates heap memory for the struct.
// 2. Creates a stack slot (alloca).
// 3. Stores the heap address into the stack slot.
func NewClass(block *bc.BlockHolder, name string, rawType types.Type) *Class {

	// 1. Normalize Type: Ensure we have a pointer-to-struct
	var ptrType *types.PointerType
	switch t := rawType.(type) {
	case *types.StructType:
		ptrType = types.NewPointer(t)
	case *types.PointerType:
		if _, ok := t.ElemType.(*types.StructType); !ok {
			panic(fmt.Sprintf("NewClass expects pointer-to-struct, got pointer to %T", t.ElemType))
		}
		ptrType = t
	default:
		panic(fmt.Sprintf("NewClass expects struct or pointer-to-struct, got %T", rawType))
	}

	// 2. Heap Allocation (malloc)
	// Calculate size: GetElementPtr hack to get size of the underlying struct
	zero := constant.NewNull(ptrType)
	one := constant.NewInt(types.I32, 1)
	gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
	size := constant.NewPtrToInt(gep, types.I64)

	// Call runtime allocator (malloc)
	// Note: Replace c.Instance.Funcs[c.ALLOC] with your specific allocator function lookup
	mallocCall := block.N.NewCall(c.Instance.Funcs[c.ALLOC], size)

	// Cast i8* (from malloc) to %MyStruct*
	heapPtr := block.N.NewBitCast(mallocCall, ptrType)

	// 3. Stack Allocation (The "Slot")
	// Create a slot on the stack that holds a %MyStruct*
	// LLVM IR: %ptr = alloca %MyStruct*
	stackSlot := block.N.NewAlloca(ptrType)

	// 4. Initialize Slot
	// Store the heap address into the stack slot
	// LLVM IR: store %MyStruct* %heapPtr, %MyStruct** %stackSlot
	block.N.NewStore(heapPtr, stackSlot)

	return &Class{
		Name: name,
		UDT:  ptrType,   // %MyStruct*
		Ptr:  stackSlot, // %MyStruct**
	}
}

// Update changes the value inside the stack slot.
// Basically: variable = newValue
func (s *Class) Update(bh *bc.BlockHolder, v value.Value) {
	block := bh.N

	if v == nil {
		errorutils.Abort(errorutils.InternalError, "cannot update object with nil value")
	}

	ptrType := s.UDT
	// Sanity check: Ensure our slot exists
	if s.Ptr == nil {
		// A. Allocate new heap memory for the struct instance
		zero := constant.NewNull(ptrType)
		one := constant.NewInt(types.I32, 1)
		gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
		size := constant.NewPtrToInt(gep, types.I64)

		// Call runtime allocator (Note: Using c.Instance.Funcs[c.ALLOC] placeholder)
		mallocCall := block.NewCall(c.Instance.Funcs[c.ALLOC], size)
		heapPtr := block.NewBitCast(mallocCall, ptrType) // %MyStruct*

		// B. Create the stack slot (alloca)
		// LLVM IR: %ptr_slot = alloca %MyStruct*
		s.Ptr = block.NewAlloca(ptrType)

		// C. Initialize slot with the new heap pointer
		// LLVM IR: store %MyStruct* %heapPtr, %MyStruct** %ptr_slot
		block.NewStore(heapPtr, s.Ptr)
	}

	// === Scenario 1: Input is a Pointer (e.g., %MyStruct*) ===
	// We simply overwrite the address stored in the slot.
	if v.Type().Equal(s.UDT) {
		// LLVM IR: store %MyStruct* %v, %MyStruct** %s.Ptr
		block.NewStore(v, s.Ptr)
		return
	}

	// === Scenario 2: Input is a Value (e.g., %MyStruct) ===
	// If the user passes a raw struct value, we can't store a struct into a pointer slot directly.
	// We must allocate new heap memory for this value, then point the slot to it.
	if v.Type().Equal(s.UDT.ElemType) {

		// 1. Allocate new heap memory
		zero := constant.NewNull(s.UDT)
		one := constant.NewInt(types.I32, 1)
		gep := constant.NewGetElementPtr(s.UDT.ElemType, zero, one)
		size := constant.NewPtrToInt(gep, types.I64)

		mem := block.NewCall(c.Instance.Funcs[c.ALLOC], size)
		newHeapPtr := block.NewBitCast(mem, s.UDT)

		// 2. Store the raw value into the new heap memory
		// LLVM IR: store %MyStruct %v, %MyStruct* %newHeapPtr
		block.NewStore(v, newHeapPtr)

		// 3. Update the slot to point to this new memory
		// LLVM IR: store %MyStruct* %newHeapPtr, %MyStruct** %s.Ptr
		block.NewStore(newHeapPtr, s.Ptr)
		return
	}

	// Error: Type mismatch
	errorutils.Abort(errorutils.InternalError,
		fmt.Sprintf("Type mismatch in Update. Expected %s or %s, got %s",
			s.UDT, s.UDT.ElemType, v.Type()))
}

// Load retrieves the current object pointer from the stack slot.
func (s *Class) Load(bh *bc.BlockHolder) value.Value {
	// LLVM IR: %val = load %MyStruct*, %MyStruct** %s.Ptr
	return bh.N.NewLoad(s.UDT, s.Ptr)
}

func (s *Class) FieldPtr(block *bc.BlockHolder, idx int) value.Value {
	zero := constant.NewInt(types.I32, 0)
	i := constant.NewInt(types.I32, int64(idx))

	// unwrap pointer-to-struct
	sPtr := s.UDT
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

	sPtr := s.UDT

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
