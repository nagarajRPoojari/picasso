package utils

import "github.com/llir/llvm/ir/types"

func GetTypeString(t types.Type) string {
	var target string
	switch et := t.(type) {
	case *types.PointerType:
		if st, ok := et.ElemType.(*types.StructType); ok {
			target = st.Name()
			if target == "" {
				target = st.String()
			}
			if target[0:1] == "%" {
				target = target[1 : len(target)-1]
			}
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
