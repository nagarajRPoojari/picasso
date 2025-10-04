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

// Array struct describes a 1D runtime Array object.
type Array struct {
	Ptr       value.Value
	ElemType  types.Type
	ArrayType *types.StructType
	Dims      []value.Value
}

func NewArray(block *ir.Block, elemType types.Type, dims []value.Value) *Array {
	// Struct layout: { i64 length, i8* data, i64* shape, i64 rank }
	arrType := types.NewStruct(
		types.I64,                   // length
		types.NewPointer(types.I8),  // data
		types.NewPointer(types.I64), // shape (i64*)
		types.I64,                   // rank
	)

	allocFn := gc.Instance.ArrayAlloc()
	if allocFn == nil {
		panic("lang_alloc_array not declared in module (gc.Instance.ArrayAlloc returned nil)")
	}

	totalLen := dims[0]
	for i := 1; i < len(dims); i++ {
		totalLen = block.NewMul(totalLen, dims[i])
	}

	elemSize := constant.NewInt(types.I64, int64(sizeOf(elemType)))

	structAlloc := block.NewCall(allocFn, totalLen, elemSize) // returns Array*

	shapeCount := constant.NewInt(types.I64, int64(len(dims)))
	shapeElemSize := constant.NewInt(types.I64, 8)

	shapeStruct := block.NewCall(allocFn, shapeCount, shapeElemSize) // returns Array* for the shape container
	shapeDataFieldPtr := block.NewGetElementPtr(arrType, shapeStruct,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	shapeRaw := block.NewLoad(types.NewPointer(types.I8), shapeDataFieldPtr) // i8* raw buffer
	shapePtrCast := block.NewBitCast(shapeRaw, types.NewPointer(types.I64))  // i64* shapePtr

	for i, d := range dims {
		elemPtr := block.NewGetElementPtr(types.I64, shapePtrCast, constant.NewInt(types.I64, int64(i)))
		block.NewStore(d, elemPtr)
	}

	rank := constant.NewInt(types.I64, int64(len(dims)))

	shapeFieldPtr := block.NewGetElementPtr(arrType, structAlloc,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	block.NewStore(shapePtrCast, shapeFieldPtr)

	rankPtr := block.NewGetElementPtr(arrType, structAlloc,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 3),
	)
	block.NewStore(rank, rankPtr)

	return &Array{
		Ptr:       structAlloc,
		ElemType:  elemType,
		ArrayType: arrType,
		Dims:      dims,
	}
}
func sizeOf(ty types.Type) int {
	switch t := ty.(type) {
	case *types.IntType:
		return int(t.BitSize / 8)
	case *types.FloatType:
		panic(fmt.Sprintf("unsupported type in sizeOf: %T", ty))
	case *types.PointerType:
		return 8 // assuming 64-bit target
	default:
		panic(fmt.Sprintf("unsupported type in sizeOf: %T", ty))
	}
}

func (a *Array) Slot() value.Value { return a.Ptr }

func (a *Array) Type() types.Type {
	return types.NewPointer(a.ArrayType)
}

func (a *Array) Cast(block *ir.Block, v value.Value) (value.Value, error) {
	if v == nil {
		return nil, errorsx.NewCompilationError("cannot cast nil to array")
	}

	targetTy := types.NewPointer(a.ArrayType)

	if vt := v.Type(); vt != nil && vt.Equal(targetTy) {
		return v, nil
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		casted := block.NewBitCast(v, targetTy)
		return casted, nil
	}

	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to array", v))
}

func (a *Array) NativeTypeString() string { return "array" }

func (a *Array) Len(block *ir.Block) value.Value {
	lengthPtr := block.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	return block.NewLoad(types.I64, lengthPtr)
}

func (a *Array) Load(block *ir.Block) value.Value {
	return block.NewLoad(a.ArrayType, a.Ptr)
}

func (a *Array) Update(block *ir.Block, v value.Value) {
	block.NewStore(a.Ptr, v)
}

// LoadRank returns i64 rank field from runtime struct
func (a *Array) LoadRank(block *ir.Block) value.Value {
	rankPtr := block.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 3),
	)
	return block.NewLoad(types.I64, rankPtr)
}

// LoadShapePtr returns i64* pointer to shape buffer
func (a *Array) LoadShapePtr(block *ir.Block) value.Value {
	shapePtrField := block.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	raw := block.NewLoad(types.NewPointer(types.I64), shapePtrField)
	return raw
}
func (a *Array) IndexOffset(block *ir.Block, indices []value.Value) value.Value {
	shapePtr := a.LoadShapePtr(block)                      // i64*
	var offset value.Value = constant.NewInt(types.I64, 0) // offset accumulator

	for i := 0; i < len(indices); i++ {
		var prod value.Value = constant.NewInt(types.I64, 1) // product of remaining dims

		for j := i + 1; j < len(indices); j++ {
			elemPtr := block.NewGetElementPtr(types.I64, shapePtr, constant.NewInt(types.I64, int64(j)))
			dimVal := block.NewLoad(types.I64, elemPtr)

			prod = block.NewMul(prod, dimVal) // *ir.InstMul satisfies value.Value
		}

		offsetPart := block.NewMul(indices[i], prod) // also value.Value
		offset = block.NewAdd(offset, offsetPart)    // *ir.InstAdd satisfies value.Value
	}

	return offset
}

func (a *Array) StoreByIndex(block *ir.Block, indices []value.Value, val value.Value) {
	offset := a.IndexOffset(block, indices)

	// Load data pointer
	dataPtrField := block.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	raw := block.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := block.NewBitCast(raw, types.NewPointer(a.ElemType))

	elemPtr := block.NewGetElementPtr(a.ElemType, elemsPtr, offset)
	block.NewStore(val, elemPtr)
}

func (a *Array) LoadByIndex(block *ir.Block, indices []value.Value) value.Value {
	offset := a.IndexOffset(block, indices)

	dataPtrField := block.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	raw := block.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := block.NewBitCast(raw, types.NewPointer(a.ElemType))

	elemPtr := block.NewGetElementPtr(a.ElemType, elemsPtr, offset)
	return block.NewLoad(a.ElemType, elemPtr)
}
