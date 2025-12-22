package typedef

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

type MethodSig struct {
	Name     string
	Hash     uint32
	FuncType *ir.Func
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
