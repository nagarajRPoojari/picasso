package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/x-lang/ast"
)

type MetaClass struct {
	VarIndexMap map[string]int
	VarAST      map[string]*ast.VariableDeclarationStatement

	Methods map[string]*ir.Func

	UDT types.Type
}

func NewMetaClass() *MetaClass {
	return &MetaClass{
		VarIndexMap: make(map[string]int),
		VarAST:      make(map[string]*ast.VariableDeclarationStatement),
		Methods:     make(map[string]*ir.Func),
	}
}
func (m *MetaClass) FieldType(idx int) types.Type {
	return m.UDT.(*types.StructType).Fields[idx]
}
