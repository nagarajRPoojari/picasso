package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
	errorutils "github.com/nagarajRPoojari/picasso/irgen/codegen/error"
)

type MetaClass struct {
	FieldIndexMap map[string]int
	VarAST        map[string]*ast.VariableDeclarationStatement

	InternalFields map[string]struct{}

	Methods    map[string]*ir.Func
	MethodArgs map[string][]ast.Type
	Returns    map[string]ast.Type

	ArrayVarsEleTypes map[int]types.Type
	// UDT is pointer-to-struct
	UDT types.Type

	Implements string

	Internal bool
}

func (mc *MetaClass) StructType() *types.StructType {
	ptr, ok := mc.UDT.(*types.PointerType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalTypeError, "UDT is not a pointer-to-struct")
	}
	st, ok := ptr.ElemType.(*types.StructType)
	if !ok {
		errorutils.Abort(errorutils.InternalError, errorutils.InternalTypeError, "UDT pointer does not point to a struct")
	}
	return st
}

func NewMetaClass(udt *types.PointerType, implements string) *MetaClass {
	return &MetaClass{
		FieldIndexMap:     make(map[string]int),
		ArrayVarsEleTypes: make(map[int]types.Type),
		InternalFields:    make(map[string]struct{}),
		VarAST:            make(map[string]*ast.VariableDeclarationStatement),
		UDT:               udt,
		Methods:           make(map[string]*ir.Func),
		MethodArgs:        make(map[string][]ast.Type),
		Returns:           map[string]ast.Type{},
		Implements:        implements,
	}
}
func (m *MetaClass) FieldType(idx int) types.Type {
	return m.StructType().Fields[idx]
}
