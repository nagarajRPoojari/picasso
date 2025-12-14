package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
	"github.com/nagarajRPoojari/x-lang/generator/c"
	"github.com/nagarajRPoojari/x-lang/generator/handlers/constants"
	bc "github.com/nagarajRPoojari/x-lang/generator/type/block"
)

// Array struct describes a 1D runtime Array object.
type Array struct {
	Ptr       value.Value
	ElemType  types.Type
	ArrayType *types.StructType
}

var ARRAYSTRUCT = types.NewStruct(
	types.NewPointer(types.I8),  // data
	types.NewPointer(types.I64), // shape (i64*)
	types.I64,                   // length
	types.I64,                   // rank
)

func init() {
	ARRAYSTRUCT.SetName(constants.ARRAY)
}

// Assuming bh, c, value, types, and constant are correctly defined/imported.
func NewArray(bh *bc.BlockHolder, elemType types.Type, eleSize value.Value, dims []value.Value) *Array {
	allocFn := c.Instance.Funcs[c.FUNC_ARRAY_ALLOC]

	totalLen := dims[0]
	for i := 1; i < len(dims); i++ {
		totalLen = bh.N.NewMul(totalLen, dims[i])
	}

	// Rank is fine as I32 for the call, as the C function expects an 'int' (typically 32-bit).
	rankVal := constant.NewInt(types.I32, int64(len(dims)))

	structAlloc := bh.N.NewCall(allocFn, totalLen, eleSize, rankVal)
	expectedPtrType := types.NewPointer(ARRAYSTRUCT)
	structPtr := bh.N.NewBitCast(structAlloc, expectedPtrType)

	shapeFieldPtr := bh.N.NewGetElementPtr(ARRAYSTRUCT, structPtr, // Use structPtr here
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1), // Index 1 is 'shape'
	)
	shapeBufPtr := bh.N.NewLoad(types.NewPointer(types.I64), shapeFieldPtr)

	for i, d := range dims {
		elemPtr := bh.N.NewGetElementPtr(types.I64, shapeBufPtr, constant.NewInt(types.I32, int64(i)))

		bh.N.NewStore(d, elemPtr)
	}

	return &Array{
		Ptr:       structAlloc,
		ElemType:  elemType,
		ArrayType: ARRAYSTRUCT,
	}
}

func (a *Array) Slot() value.Value { return a.Ptr }

func (a *Array) Type() types.Type {
	return types.NewPointer(a.ArrayType)
}

func (a *Array) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	if v == nil {
		return nil, errorsx.NewCompilationError("cannot cast nil to array")
	}

	targetTy := types.NewPointer(a.ArrayType)

	if vt := v.Type(); vt != nil && vt.Equal(targetTy) {
		return v, nil
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		casted := block.N.NewBitCast(v, targetTy)
		return casted, nil
	}

	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to array", v))
}

func (a *Array) NativeTypeString() string { return "array" }

func (a *Array) Len(block *bc.BlockHolder) value.Value {
	lengthPtr := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	return block.N.NewLoad(types.I64, lengthPtr)
}

func (a *Array) Load(block *bc.BlockHolder) value.Value {
	return a.Ptr
}

func (a *Array) Update(block *bc.BlockHolder, v value.Value) {
	block.N.NewStore(a.Ptr, v)
}

func (a *Array) UpdateV2(block *bc.BlockHolder, v *Array) {
	*a = *v
}

// LoadRank returns i64 rank field from runtime struct
func (a *Array) LoadRank(block *bc.BlockHolder) value.Value {
	rankPtr := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 3),
	)
	return block.N.NewLoad(types.I64, rankPtr)
}

// LoadShapePtr returns i64* pointer to shape buffer
func (a *Array) LoadShapePtr(block *bc.BlockHolder) value.Value {
	shapePtrField := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	raw := block.N.NewLoad(types.NewPointer(types.I64), shapePtrField)
	return raw
}

func (a *Array) IndexOffset(block *bc.BlockHolder, indices []value.Value) value.Value {
	shapePtr := a.LoadShapePtr(block)
	var offset value.Value = constant.NewInt(types.I64, 0)

	for i := range indices {
		var prod value.Value = constant.NewInt(types.I64, 1)

		for j := i + 1; j < len(indices); j++ {
			elemPtr := block.N.NewGetElementPtr(types.I64, shapePtr, constant.NewInt(types.I64, int64(j)))
			dimVal := block.N.NewLoad(types.I64, elemPtr)

			prod = block.N.NewMul(prod, dimVal)
		}

		offsetPart := block.N.NewMul(indices[i], prod)
		offset = block.N.NewAdd(offset, offsetPart)
	}

	return offset
}

// StoreByIndex updates element value at given index
func (a *Array) StoreByIndex(block *bc.BlockHolder, indices []value.Value, val value.Value) {
	offset := a.IndexOffset(block, indices)
	dataPtrField := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	raw := block.N.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := block.N.NewBitCast(raw, types.NewPointer(a.ElemType))

	elemPtr := block.N.NewGetElementPtr(a.ElemType, elemsPtr, offset)
	block.N.NewStore(val, elemPtr)
}

// LoadByIndex retrieves element value at given index
func (a *Array) LoadByIndex(block *bc.BlockHolder, indices []value.Value) value.Value {
	offset := a.IndexOffset(block, indices)

	dataPtrField := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	raw := block.N.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := block.N.NewBitCast(raw, types.NewPointer(a.ElemType))

	elemPtr := block.N.NewGetElementPtr(a.ElemType, elemsPtr, offset)
	return block.N.NewLoad(a.ElemType, elemPtr)
}
