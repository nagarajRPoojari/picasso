package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/c"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
	"github.com/nagarajRPoojari/niyama/irgen/codegen/handlers/utils"
	bc "github.com/nagarajRPoojari/niyama/irgen/codegen/type/block"
)

type Class struct {
	Name string             // Name for debugging/lookup
	UDT  *types.PointerType // The type of the object (e.g., %MyStruct*)
	Ptr  value.Value        // The stack slot (alloca) holding the pointer (e.g., %MyStruct**)
}

// NewClass creates a new object instance by allocation memory in heap.
// Key Logic:
//   - Allocate memory in heap & store its pointer in stack slot.
//   - Memory not deallocated even if pointer in stack slot goes out of scope.
func NewClass(block *bc.BlockHolder, name string, rawType types.Type) *Class {

	// normalize Type: Ensure we have a pointer-to-struct
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

	// Calculate size: GetElementPtr hack to get size of the underlying struct
	zero := constant.NewNull(ptrType)
	one := constant.NewInt(types.I32, 1)
	gep := constant.NewGetElementPtr(ptrType.ElemType, zero, one)
	size := constant.NewPtrToInt(gep, types.I64)

	mallocCall := block.N.NewCall(c.Instance.Funcs[c.FUNC_ALLOC], size)
	heapPtr := block.N.NewBitCast(mallocCall, ptrType)

	// Create a slot on the stack that holds a %MyStruct*
	// LLVM IR: %ptr = alloca %MyStruct*
	stackSlot := block.N.NewAlloca(ptrType)

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
	if s.Ptr == nil {
		// cases do exist where we need to create Class instance with memory allocated beforehand.
		// example in assignment, where i avoid NewClass to avoid heap allocation, in such cases s.Ptr
		// remains null, so allocate stack slot.
		s.Ptr = block.NewAlloca(ptrType)
	}

	// We simply overwrite the address stored in the slot.
	if v.Type().Equal(s.UDT) {
		block.NewStore(v, s.Ptr)
		return
	}

	// If the user passes a raw struct value, we can't store a struct into a pointer slot directly.
	// I must allocate new heap memory for this value, then point the slot to it.
	if v.Type().Equal(s.UDT.ElemType) {
		block.NewStore(v, s.Ptr)
		return
	}

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

func (s *Class) UpdateField(bh *bc.BlockHolder, th *TypeHandler, idx int, v value.Value, expected types.Type) {
	block := bh.N
	if v == nil {
		// errorsx.PanicCompilationError(fmt.Sprintf("cannot update field with nil value: %v", v))
		return
	}
	val := ensureType(bh, th, v, expected)
	fieldPtr := s.FieldPtr(bh, idx)
	block.NewStore(val, fieldPtr)
}

func (s *Class) LoadField(block *bc.BlockHolder, idx int, fieldType types.Type) value.Value {
	fieldPtr := s.FieldPtr(block, idx)

	return block.N.NewLoad(fieldType, fieldPtr)
}

func (s *Class) Cast(bh *bc.BlockHolder, v value.Value) (value.Value, error) {
	return nil, nil
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

// heuristic type check for Class types.
func ensureType(block *bc.BlockHolder, th *TypeHandler, v value.Value, target types.Type) value.Value {
	var ret value.Value
	var err error

	ret, err = ensureClassType(block, th, v, target)

	if err != nil {
		ret, err = ensureInterfaceType(block, th, v, target)
		if err != nil {
			panic(err)
		}
	}

	return ret
}

// heuristic type check for Class types.
func ensureClassType(block *bc.BlockHolder, _ *TypeHandler, v value.Value, target types.Type) (value.Value, error) {
	srcType := v.Type()

	if srcType.Equal(target) {
		return v, nil
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
				return block.N.NewBitCast(v, target), nil
			}

			// Child -> Parent layout check
			if len(srcStruct.Fields) < len(tgtStruct.Fields) {
				return nil, fmt.Errorf("ensureType: source struct has fewer fields than target parent")
			}

			for i := range tgtStruct.Fields {
				if !srcStruct.Fields[i].Equal(tgtStruct.Fields[i]) {
					return nil, fmt.Errorf("ensureType: field %d type mismatch: %v vs %v", i, srcStruct.Fields[i], tgtStruct.Fields[i])
				}
			}

			// Safe to bitcast child pointer to parent pointer
			return block.N.NewBitCast(v, target), nil
		}

		// Not structs, fallback: generic pointer bitcast
		return block.N.NewBitCast(v, target), nil
	}

	// Bitcast function pointers
	if _, okSrcFunc := srcType.(*types.FuncType); okSrcFunc {
		if _, okTgtFunc := target.(*types.FuncType); okTgtFunc {
			return block.N.NewBitCast(v, target), nil
		}
	}
	return nil, fmt.Errorf("ensureType: cannot convert %v -> %v", srcType, target)
}

func ensureInterfaceType(block *bc.BlockHolder, th *TypeHandler, v value.Value, target types.Type) (value.Value, error) {

	for _, cls := range th.InterfaceUDTS[utils.GetTypeString(target)].ImplementedBy {
		clsUDT := th.ClassUDTS[cls]
		if ret, err := ensureClassType(block, th, v, clsUDT.UDT); err == nil {
			return ret, nil
		}
	}
	return nil, fmt.Errorf("ensureType: cannot convert %v -> %v", v.Type(), target)
}

func (f *Class) NativeTypeString() string { return f.Name }
