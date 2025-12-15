package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/niyama/irgen/ast"
	errorutils "github.com/nagarajRPoojari/niyama/irgen/codegen/error"
)

type MetaClass struct {
	FieldIndexMap map[string]int
	VarAST        map[string]*ast.VariableDeclarationStatement

	Methods map[string]*ir.Func
	Returns map[string]ast.Type

	ArrayVarsEleTypes map[int]types.Type
	// UDT is pointer-to-struct
	UDT types.Type
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

func NewMetaClass() *MetaClass {
	return &MetaClass{
		FieldIndexMap:     make(map[string]int),
		ArrayVarsEleTypes: make(map[int]types.Type),
		VarAST:            make(map[string]*ast.VariableDeclarationStatement),
		Methods:           make(map[string]*ir.Func),
		Returns:           map[string]ast.Type{},
	}
}
func (m *MetaClass) FieldType(idx int) types.Type {
	return m.StructType().Fields[idx]
}
