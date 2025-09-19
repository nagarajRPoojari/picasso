package compiler

import "github.com/llir/llvm/ir/types"

func GetTypeString(t types.Type) string {
	var target string
	switch et := t.(type) {
	case *types.PointerType:
		if st, ok := et.ElemType.(*types.StructType); ok {
			target = st.Name()
		} else {
			target = t.String()
		}
	case *types.StructType:
		target = et.Name()
	default:
		target = t.String()
	}
	return target
}

type IdentifierBuilder struct {
	module string
}

func NewIdentifierBuilder(module string) *IdentifierBuilder {
	return &IdentifierBuilder{module}
}

func (t *IdentifierBuilder) Attach(name ...string) string {
	res := t.module
	for _, n := range name {
		res += "." + n
	}
	return res
}
