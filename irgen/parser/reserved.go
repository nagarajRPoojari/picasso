package parser

var reserved_keywords map[string]struct{}

const (
	THREAD = "thread"
	FUNC   = "fn"
	VOID   = "void"
	BOOL   = "bool"
	CHAR   = "char"
	STRING = "string"

	INT8 = "int8"
	I8   = "i8"

	INT16 = "int16"
	I16   = "i16"

	INT32 = "int32"
	I32   = "i32"

	INT64 = "int64"
	I64   = "i64"

	INT = "int"

	UINT8 = "uint8"
	U8    = "u8"

	UINT16 = "uint16"
	U16    = "u16"

	UINT32 = "uint32"
	U32    = "u32"

	UINT64 = "uint64"
	U64    = "u64"

	UINT = "uint"

	FLOAT16 = "float16"
	F16     = "f16"

	FLOAT32 = "float32"
	F32     = "f32"

	FLOAT64 = "float64"
	F64     = "f64"

	ATOMIC_INT8  = "atomic_int8"
	ATOMIC_INT16 = "atomic_int16"
	ATOMIC_INT32 = "atomic_int32"
	ATOMIC_INT64 = "atomic_int64"

	ATOMIC_UINT8  = "atomic_uint8"
	ATOMIC_UINT16 = "atomic_uint16"
	ATOMIC_UINT32 = "atomic_uint32"
	ATOMIC_UINT64 = "atomic_uint64"

	PTR    = "ptr"
	ARRAY  = "array"
	SLICE  = "slice"
	MAP    = "map"
	STRUCT = "struct"
	UNION  = "union"
	ENUM   = "enum"
)

func init() {
	reserved_keywords = make(map[string]struct{})

	keywords := []string{
		/* Core / Runtime Types */
		THREAD,
		FUNC,
		VOID,
		BOOL,
		CHAR,
		STRING,

		/* Signed Integers */
		INT8, I8,
		INT16, I16,
		INT32, I32,
		INT64, I64,
		INT,

		/* Unsigned Integers */
		UINT8, U8,
		UINT16, U16,
		UINT32, U32,
		UINT64, U64,
		UINT,

		/* Floating Point */
		FLOAT16, F16,
		FLOAT32, F32,
		FLOAT64, F64,

		/* Atomics */
		ATOMIC_INT8,
		ATOMIC_INT16,
		ATOMIC_INT32,
		ATOMIC_INT64,
		ATOMIC_UINT8,
		ATOMIC_UINT16,
		ATOMIC_UINT32,
		ATOMIC_UINT64,

		/* Pointer / Composite */
		PTR,
		ARRAY,
		SLICE,
		MAP,
		STRUCT,
		UNION,
		ENUM,
	}

	for _, kw := range keywords {
		reserved_keywords[kw] = struct{}{}
	}
}
