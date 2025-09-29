package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
	errorsx "github.com/nagarajRPoojari/x-lang/error"
)

type MetaClass struct {
	FieldIndexMap map[string]int
	VarAST        map[string]*ast.VariableDeclarationStatement

	Methods map[string]*ir.Func

	// UDT is pointer-to-struct
	UDT types.Type
}

func (mc *MetaClass) StructType() *types.StructType {
	ptr, ok := mc.UDT.(*types.PointerType)
	if !ok {
		errorsx.PanicCompilationError("UDT is not a pointer-to-struct")
	}
	st, ok := ptr.ElemType.(*types.StructType)
	if !ok {
		errorsx.PanicCompilationError("UDT pointer does not point to a struct")
	}
	return st
}

func NewMetaClass() *MetaClass {
	return &MetaClass{
		FieldIndexMap: make(map[string]int),
		VarAST:        make(map[string]*ast.VariableDeclarationStatement),
		Methods:       make(map[string]*ir.Func),
	}
}
func (m *MetaClass) FieldType(idx int) types.Type {
	return m.StructType().Fields[idx]
}
