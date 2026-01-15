package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/c"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/constants"
	rterr "github.com/nagarajRPoojari/picasso/irgen/codegen/libs/private/runtime"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/type/primitives/ints"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
)

// Array struct describes a 1D runtime Array object.
type Array struct {
	Ptr       value.Value
	ElemType  types.Type
	ArrayType *types.StructType

	ElementTypeString string
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
func NewArray(bh *bc.BlockHolder, elemType types.Type, eleSize value.Value, dims []value.Value, ElementTypeString string) *Array {
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
		Ptr:               structAlloc,
		ElemType:          elemType,
		ArrayType:         ARRAYSTRUCT,
		ElementTypeString: ElementTypeString,
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

func (a *Array) Len(block *bc.BlockHolder) *ints.Int64 {
	lengthPtr := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	return &ints.Int64{
		NativeType: types.I64,
		Value:      lengthPtr,
	}
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
func (a *Array) LoadRank(block *bc.BlockHolder) *ints.Int64 {
	rankPtr := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 3),
	)
	return &ints.Int64{
		NativeType: types.I64,
		Value:      rankPtr,
	}
}

// LoadShapePtr returns i64* pointer to shape buffer
func (a *Array) LoadShapePtr(block *bc.BlockHolder) value.Value {
	shapePtrField := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)
	raw := block.N.NewLoad(types.NewPointer(types.I64), shapePtrField)
	return raw
}

func (a *Array) LoadShapeArray(block *bc.BlockHolder) *Array {
	b := block.N

	rank := a.LoadRank(block).Load(block)
	elemSize := constant.NewInt(types.I64, 8)

	rankVal := constant.NewInt(types.I32, 1)

	structAlloc := b.NewCall(c.Instance.Funcs[c.FUNC_ARRAY_ALLOC], rank, elemSize, rankVal)
	arrayPtr := b.NewBitCast(structAlloc, types.NewPointer(ARRAYSTRUCT))

	origShapePtr := a.LoadShapePtr(block)

	dataField := b.NewGetElementPtr(
		ARRAYSTRUCT,
		arrayPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)

	dataAsI8 := b.NewBitCast(origShapePtr, types.NewPointer(types.I8))
	b.NewStore(dataAsI8, dataField)

	shapeField := b.NewGetElementPtr(
		ARRAYSTRUCT,
		arrayPtr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 1),
	)

	shapeBuf := b.NewLoad(types.NewPointer(types.I64), shapeField)

	firstDim := b.NewGetElementPtr(
		types.I64,
		shapeBuf,
		constant.NewInt(types.I32, 0),
	)
	b.NewStore(rank, firstDim)

	return &Array{
		Ptr:               arrayPtr,
		ElemType:          types.I64,
		ArrayType:         ARRAYSTRUCT,
		ElementTypeString: "int64",
	}
}

// IndexOffset computes linear offset with bounds checks
func (a *Array) IndexOffset(block *bc.BlockHolder, indices []value.Value) value.Value {
	shapePtr := a.LoadShapePtr(block)
	var offset value.Value = constant.NewInt(types.I64, 0)

	for i := range indices {
		idx := indices[i]

		// index >= 0
		checkIntCond(block, idx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")

		// index < shape[i]
		shapeElemPtr := block.N.NewGetElementPtr(types.I64, shapePtr, constant.NewInt(types.I64, int64(i)))
		dimVal := block.N.NewLoad(types.I64, shapeElemPtr)

		checkIntCond(block, idx, dimVal, enum.IPredSLT, "array index out of bounds\n")

		var prod value.Value = constant.NewInt(types.I64, 1)
		for j := i + 1; j < len(indices); j++ {
			nextDimPtr := block.N.NewGetElementPtr(types.I64, shapePtr, constant.NewInt(types.I64, int64(j)))
			nextDim := block.N.NewLoad(types.I64, nextDimPtr)
			prod = block.N.NewMul(prod, nextDim)
		}

		offset = block.N.NewAdd(offset, block.N.NewMul(idx, prod))
	}

	return offset
}

// StoreByIndex updates element value at given index
func (a *Array) StoreByIndex(block *bc.BlockHolder, indices []value.Value, val value.Value) {
	offset := a.IndexOffset(block, indices)
	dataPtrField := block.N.NewGetElementPtr(a.ArrayType, a.Ptr,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
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
		constant.NewInt(types.I32, 0),
	)
	raw := block.N.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := block.N.NewBitCast(raw, types.NewPointer(a.ElemType))

	elemPtr := block.N.NewGetElementPtr(a.ElemType, elemsPtr, offset)
	return block.N.NewLoad(a.ElemType, elemPtr)
}

func checkIntCond(block *bc.BlockHolder, v1, v2 value.Value, pred enum.IPred, errMsg string) {
	b := block.N
	passBlk := b.Parent.NewBlock("")
	failBlk := b.Parent.NewBlock("")
	cond := b.NewICmp(pred, v1, v2)
	b.NewCondBr(cond, passBlk, failBlk)
	rterr.Instance.RaiseRTError(failBlk, errMsg)
	block.Update(block.V, passBlk)
}
