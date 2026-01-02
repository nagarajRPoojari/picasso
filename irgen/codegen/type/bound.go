package typedef

import "github.com/llir/llvm/ir/types"

var intMax = map[*types.IntType]int64{
	types.I1:  1,
	types.I8:  127,
	types.I16: 32767,
	types.I32: 2147483647,
	types.I64: 9223372036854775807,
}

var intMin = map[*types.IntType]int64{
	types.I1:  0,
	types.I8:  -128,
	types.I16: -32768,
	types.I32: -2147483648,
	types.I64: -9223372036854775808,
}

var uintMax = map[*types.IntType]uint64{
	types.I1:  1,                    // 2^1  - 1
	types.I8:  255,                  // 2^8  - 1
	types.I16: 65535,                // 2^16 - 1
	types.I32: 4294967295,           // 2^32 - 1
	types.I64: 18446744073709551615, // 2^64 - 1
}

var uintMin = map[*types.IntType]uint64{
	types.I1:  0,
	types.I8:  0,
	types.I16: 0,
	types.I32: 0,
	types.I64: 0,
}

var floatMax = map[*types.FloatType]float64{
	types.Half:   65504.0,
	types.Float:  3.4028235e38,
	types.Double: 1.7976931348623157e308,
}

var floatMin = map[*types.FloatType]float64{
	types.Half:   -65504.0,
	types.Float:  -3.4028235e38,
	types.Double: -1.7976931348623157e308,
}
