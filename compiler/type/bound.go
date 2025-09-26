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

var floatMax = map[*types.FloatType]float64{
	types.Half:   65504.0, // will be converted to half
	types.Float:  3.4028235e38,
	types.Double: 1.7976931348623157e308,
}

var floatMin = map[*types.FloatType]float64{
	types.Half:   -65504.0,
	types.Float:  -3.4028235e38,
	types.Double: -1.7976931348623157e308,
}
