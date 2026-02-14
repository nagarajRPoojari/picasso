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

type Array struct {
	Ptr       value.Value
	ElemType  types.Type
	ArrayType *types.StructType

	ElementTypeString string
	Rank              int // Number of dimensions (compile-time known)
}

var ARRAYSTRUCT = types.NewStruct(
	types.NewPointer(types.I8),  // data
	types.NewPointer(types.I64), // shape (i64*)
	types.I64,                   // length
	types.I64,                   // rank
	types.I64,                   // capacity
	types.I32,                   // elesize
)

func init() {
	ARRAYSTRUCT.SetName(constants.ARRAY)
}

// Assuming bh, c, value, types, and constant are correctly defined/imported.
func NewArray(bh *bc.BlockHolder, elemType types.Type, eleSize value.Value, dims []value.Value, ElementTypeString string) *Array {
	allocFn := c.Instance.Funcs[c.FUNC_ARRAY_ALLOC]
	rankVal := constant.NewInt(types.I32, int64(len(dims)))

	args := []value.Value{eleSize, rankVal}
	args = append(args, dims...)
	structAlloc := bh.N.NewCall(allocFn, args...)

	return &Array{
		Ptr:               structAlloc,
		ElemType:          elemType,
		ArrayType:         ARRAYSTRUCT,
		ElementTypeString: ElementTypeString,
		Rank:              len(dims),
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
	elemSize := constant.NewInt(types.I32, 8)

	rankVal := constant.NewInt(types.I32, 1)

	structAlloc := b.NewCall(c.Instance.Funcs[c.FUNC_ARRAY_ALLOC], elemSize, rankVal, rank)
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
		Rank:              1,
	}
}

// StoreByIndex updates element value at given index in a jagged array
func (a *Array) StoreByIndex(block *bc.BlockHolder, indices []value.Value, val value.Value) {
	b := block.N
	currentArray := a.Ptr

	for i := 0; i < len(indices)-1; i++ {
		idx := indices[i]

		checkIntCond(block, idx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
		b = block.N

		lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 2),
		)
		length := b.NewLoad(types.I64, lengthPtr)
		checkIntCond(block, idx, length, enum.IPredSLT, "array index out of bounds\n")
		b = block.N

		getSubarrayFn := c.Instance.Funcs[c.FUNC_GET_SUBARRAY]
		currentArray = b.NewCall(getSubarrayFn, currentArray, idx)
	}

	lastIdx := indices[len(indices)-1]

	checkIntCond(block, lastIdx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
	b = block.N
	lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	length := b.NewLoad(types.I64, lengthPtr)
	checkIntCond(block, lastIdx, length, enum.IPredSLT, "array index out of bounds\n")
	b = block.N

	dataPtrField := b.NewGetElementPtr(a.ArrayType, currentArray,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	raw := b.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := b.NewBitCast(raw, types.NewPointer(a.ElemType))
	elemPtr := b.NewGetElementPtr(a.ElemType, elemsPtr, lastIdx)
	b.NewStore(val, elemPtr)
}

// StoreSubarrayByIndex stores a subarray at given index (partial indexing)
func (a *Array) StoreSubarrayByIndex(block *bc.BlockHolder, indices []value.Value, subarray *Array) {
	b := block.N
	currentArray := a.Ptr

	for i := 0; i < len(indices)-1; i++ {
		idx := indices[i]

		checkIntCond(block, idx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
		b = block.N

		lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 2),
		)
		length := b.NewLoad(types.I64, lengthPtr)
		checkIntCond(block, idx, length, enum.IPredSLT, "array index out of bounds\n")
		b = block.N

		getSubarrayFn := c.Instance.Funcs[c.FUNC_GET_SUBARRAY]
		currentArray = b.NewCall(getSubarrayFn, currentArray, idx)
	}

	lastIdx := indices[len(indices)-1]

	checkIntCond(block, lastIdx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
	b = block.N
	lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	length := b.NewLoad(types.I64, lengthPtr)
	checkIntCond(block, lastIdx, length, enum.IPredSLT, "array index out of bounds\n")
	b = block.N

	// Store the subarray using set_subarray
	setSubarrayFn := c.Instance.Funcs[c.FUNC_SET_SUBARRAY]
	b.NewCall(setSubarrayFn, currentArray, lastIdx, subarray.Ptr)
}

// LoadByIndex retrieves element value at given index in a jagged array
func (a *Array) LoadByIndex(block *bc.BlockHolder, indices []value.Value) value.Value {
	b := block.N
	currentArray := a.Ptr

	for i := 0; i < len(indices)-1; i++ {
		idx := indices[i]

		checkIntCond(block, idx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
		b = block.N

		lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 2),
		)
		length := b.NewLoad(types.I64, lengthPtr)
		checkIntCond(block, idx, length, enum.IPredSLT, "array index out of bounds\n")
		b = block.N

		getSubarrayFn := c.Instance.Funcs[c.FUNC_GET_SUBARRAY]
		currentArray = b.NewCall(getSubarrayFn, currentArray, idx)
	}

	lastIdx := indices[len(indices)-1]

	checkIntCond(block, lastIdx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
	b = block.N
	lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 2),
	)
	length := b.NewLoad(types.I64, lengthPtr)
	checkIntCond(block, lastIdx, length, enum.IPredSLT, "array index out of bounds\n")
	b = block.N

	dataPtrField := b.NewGetElementPtr(a.ArrayType, currentArray,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, 0),
	)
	raw := b.NewLoad(types.NewPointer(types.I8), dataPtrField)
	elemsPtr := b.NewBitCast(raw, types.NewPointer(a.ElemType))
	elemPtr := b.NewGetElementPtr(a.ElemType, elemsPtr, lastIdx)
	return b.NewLoad(a.ElemType, elemPtr)
}

// LoadSubarrayByIndex retrieves a subarray at given index (partial indexing)
func (a *Array) LoadSubarrayByIndex(block *bc.BlockHolder, indices []value.Value) *Array {
	b := block.N
	currentArray := a.Ptr

	for i := 0; i < len(indices); i++ {
		idx := indices[i]

		checkIntCond(block, idx, constant.NewInt(types.I64, 0), enum.IPredSGE, "array index < 0")
		b = block.N

		lengthPtr := b.NewGetElementPtr(a.ArrayType, currentArray,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, 2),
		)
		length := b.NewLoad(types.I64, lengthPtr)
		checkIntCond(block, idx, length, enum.IPredSLT, "array index out of bounds\n")
		b = block.N

		getSubarrayFn := c.Instance.Funcs[c.FUNC_GET_SUBARRAY]
		currentArray = b.NewCall(getSubarrayFn, currentArray, idx)
	}

	return &Array{
		Ptr:               currentArray,
		ElemType:          a.ElemType,
		ArrayType:         a.ArrayType,
		ElementTypeString: a.ElementTypeString,
		Rank:              a.Rank - len(indices),
	}
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
