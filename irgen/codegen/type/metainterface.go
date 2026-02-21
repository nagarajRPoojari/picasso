package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/nagarajRPoojari/picasso/irgen/ast"
)

type MethodSig struct {
	Name       string
	Hash       uint32
	FuncType   *ir.Func
	Parameters []ast.Parameter
	ReturnType ast.Type
}

type MetaInterface struct {
	Methods map[string]MethodSig
	UDT     types.Type

	ImplementedBy []string
}

func NewMetaInterface() *MetaInterface {
	return &MetaInterface{
		Methods:       make(map[string]MethodSig),
		ImplementedBy: make([]string, 0),
	}
}
