package utils

import (
	"fmt"
	"math"
	"reflect"
	"sort"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

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

func IsNullConstant(v value.Value) bool {
	switch v.(type) {
	case *constant.Null:
		return true
	default:
		return false
	}
}

func HashFuncSig(params, ret any) uint32 {
	h := uint32(2166136261)
	hashValue(&h, reflect.ValueOf(params))
	hashValue(&h, reflect.ValueOf(ret))
	return h
}

const fnv32Prime = 16777619

func fnv(h *uint32, b byte) {
	*h ^= uint32(b)
	*h *= fnv32Prime
}

func fnvUint64(h *uint32, x uint64) {
	for i := 0; i < 8; i++ {
		fnv(h, byte(x))
		x >>= 8
	}
}

func hashValue(h *uint32, v reflect.Value) {
	if !v.IsValid() {
		fnv(h, 0)
		return
	}

	// include kind
	fnv(h, byte(v.Kind()))

	switch v.Kind() {

	case reflect.Bool:
		if v.Bool() {
			fnv(h, 1)
		} else {
			fnv(h, 0)
		}

	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		fnvUint64(h, uint64(v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		fnvUint64(h, v.Uint())

	case reflect.Float32, reflect.Float64:
		fnvUint64(h, math.Float64bits(v.Convert(reflect.TypeOf(float64(0))).Float()))

	case reflect.String:
		s := v.String()
		fnvUint64(h, uint64(len(s)))
		for i := 0; i < len(s); i++ {
			fnv(h, s[i])
		}

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			hashValue(h, v.Field(i))
		}

	case reflect.Array, reflect.Slice:
		l := v.Len()
		fnvUint64(h, uint64(l))
		for i := 0; i < l; i++ {
			hashValue(h, v.Index(i))
		}

	case reflect.Map:
		keys := v.MapKeys()
		sort.Slice(keys, func(i, j int) bool {
			return fmt.Sprint(keys[i].Interface()) <
				fmt.Sprint(keys[j].Interface())
		})

		fnvUint64(h, uint64(len(keys)))
		for _, k := range keys {
			hashValue(h, k)
			hashValue(h, v.MapIndex(k))
		}

	case reflect.Pointer, reflect.Interface:
		if v.IsNil() {
			fnv(h, 0)
		} else {
			fnv(h, 1)
			hashValue(h, v.Elem())
		}

	default:
		panic("unsupported kind: " + v.Kind().String())
	}
}
