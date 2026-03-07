package typedef

import (
	"fmt"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	"github.com/nagarajRPoojari/picasso/irgen/codegen/handlers/utils"
	bc "github.com/nagarajRPoojari/picasso/irgen/codegen/type/block"
	errorsx "github.com/nagarajRPoojari/picasso/irgen/error"
)

// Tuple represents a composite type that combines multiple types.
// It's implemented as a C struct containing all the component types.
// This allows functions to return multiple values by returning a struct.
type Tuple struct {
	NativeType *types.StructType // The LLVM struct type
	Value      value.Value       // Pointer to the allocated struct
	Types      []Var             // The component variables/types
	TypeNames  []string          // Names of the component types
}

// NewTuple creates a new tuple with the given component types and registered struct type
// The structType should be the one registered in GlobalTypeList to ensure proper type definition
func NewTuple(block *bc.BlockHolder, structType *types.StructType, componentTypes []Var, typeNames []string) *Tuple {
	// @imp: allocate struct in stack
	slot := block.N.NewAlloca(structType)

	for i, v := range componentTypes {
		fieldPtr := block.N.NewGetElementPtr(structType, slot,
			constant.NewInt(types.I32, 0),
			constant.NewInt(types.I32, int64(i)),
		)
		block.N.NewStore(v.Load(block), fieldPtr)
	}

	return &Tuple{
		NativeType: structType,
		Value:      slot,
		Types:      componentTypes,
		TypeNames:  typeNames,
	}
}

// NewTupleFromStruct creates a tuple from an existing struct value
// The structVal is already a struct value returned from a function call
func NewTupleFromStruct(block *bc.BlockHolder, structVal value.Value, structType *types.StructType, typeNames []string) *Tuple {
	// Don't allocate - the struct value is already in a register
	// We'll extract fields directly from it
	return &Tuple{
		NativeType: structType,
		Value:      structVal, // Use the struct value directly
		TypeNames:  typeNames,
	}
}

// GetField returns the value of the field at the given index
func (t *Tuple) GetField(block *bc.BlockHolder, index int) value.Value {
	if index < 0 || index >= len(t.NativeType.Fields) {
		panic(fmt.Sprintf("tuple index %d out of bounds (size: %d)", index, len(t.NativeType.Fields)))
	}

	// Extract field from struct value using extractvalue instruction
	return block.N.NewExtractValue(t.Value, uint64(index))
}

// SetField sets the value of the field at the given index
func (t *Tuple) SetField(block *bc.BlockHolder, index int, val value.Value) {
	if index < 0 || index >= len(t.NativeType.Fields) {
		panic(fmt.Sprintf("tuple index %d out of bounds (size: %d)", index, len(t.NativeType.Fields)))
	}

	fieldPtr := block.N.NewGetElementPtr(t.NativeType, t.Value,
		constant.NewInt(types.I32, 0),
		constant.NewInt(types.I32, int64(index)),
	)
	block.N.NewStore(val, fieldPtr)
}

// Load returns the tuple struct value
func (t *Tuple) Load(block *bc.BlockHolder) value.Value {
	return block.N.NewLoad(t.NativeType, t.Value)
}

// Update stores a new tuple value
func (t *Tuple) Update(block *bc.BlockHolder, v value.Value) {
	block.N.NewStore(v, t.Value)
}

// Slot returns the pointer to the tuple
func (t *Tuple) Slot() value.Value {
	return t.Value
}

// Type returns the LLVM type of the tuple
func (t *Tuple) Type() types.Type {
	return types.NewPointer(t.NativeType)
}

// GetOrCreateTupleType is a utility function that creates or retrieves a tuple type.
func GetOrCreateTupleType(tupleType *ast.TupleType, module *ir.Module, typeHandler *TypeHandler, resolveAlias func(string) string, globalTypeList map[string]types.Type) (types.Type, []string) {
	tupleFieldTypes := make([]types.Type, len(tupleType.Types))
	typeNames := make([]string, len(tupleType.Types))

	for i, componentType := range tupleType.Types {
		resolvedType := resolveAlias(componentType.Get())
		tupleFieldTypes[i] = typeHandler.GetLLVMType(resolvedType)
		typeNames[i] = resolvedType
	}

	tupleName := utils.GenerateTupleName(tupleFieldTypes, typeNames)

	if existingType, ok := globalTypeList[tupleName]; ok {
		return existingType, typeNames
	}

	structType := types.NewStruct(tupleFieldTypes...)
	retType := module.NewTypeDef(tupleName, structType)
	globalTypeList[tupleName] = retType

	return retType, typeNames
}

// Cast attempts to cast a value to this tuple type
func (t *Tuple) Cast(block *bc.BlockHolder, v value.Value) (value.Value, error) {
	if v == nil {
		return nil, errorsx.NewCompilationError("cannot cast nil to tuple")
	}

	targetTy := types.NewPointer(t.NativeType)

	if vt := v.Type(); vt != nil && vt.Equal(targetTy) {
		return v, nil
	}

	if _, ok := v.Type().(*types.PointerType); ok {
		casted := block.N.NewBitCast(v, targetTy)
		return casted, nil
	}

	return nil, errorsx.NewCompilationError(fmt.Sprintf("failed to typecast %v to tuple", v))
}

// NativeTypeString returns the string representation of the type
func (t *Tuple) NativeTypeString() string {
	return "tuple"
}

// Constant returns a null constant for the tuple type
func (t *Tuple) Constant() constant.Constant {
	return constant.NewNull(types.NewPointer(t.NativeType))
}
