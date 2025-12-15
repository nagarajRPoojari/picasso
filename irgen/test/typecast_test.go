package test

import (
	"testing"

	"github.com/nagarajRPoojari/niyama/irgen/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBooleanTypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// boolean -> boolean
		{
			name: "boolean->boolean",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: boolean = b;
				io.printf("value=%d, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		// boolean -> int8
		{
			name: "boolean->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: int8 = b;
				io.printf("value=%d, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1, type=int8, size=1",
		},
		// boolean -> int32
		{
			name: "boolean->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: int32 = b;
				io.printf("value=%d, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1, type=int32, size=4",
		},
		// boolean -> int64
		{
			name: "boolean->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: int = b;
				io.printf("value=%d, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1, type=int64, size=8",
		},
		// boolean -> float16
		{
			name: "boolean->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: float16 = b;
				io.printf("value=%f, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1.000000, type=float16, size=4",
		},
		// boolean -> float32
		{
			name: "boolean->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: float32 = b;
				io.printf("value=%f, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1.000000, type=float32, size=4",
		},
		// boolean -> float64
		{
			name: "boolean->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say b: boolean = 1;
				say c: float64 = b;
				io.printf("value=%f, type=%s, size=%d", c, types.type(c), types.size(c));
				return 0;
			}`,
			wantOut: "value=1.000000, type=float64, size=8",
		},
		{
			name: "boolean->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say b: boolean = true;
				say s: string = b;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true}, // not allowed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestInt8TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// int8 -> boolean
		{
			name: "int8->boolean",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 1;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "int8->boolean (other than 0/1)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 100;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		// int8 -> int8
		{
			name: "int8->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		// int8 -> int32
		{
			name: "int8->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		// int8 -> int64
		{
			name: "int8->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// int8 -> float16
		{
			name: "int8->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float16, size=4",
		},
		// int8 -> float32
		{
			name: "int8->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float32, size=4",
		},
		// int8 -> float64
		{
			name: "int8->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float64, size=8",
		},
		// int8 -> string
		{
			name: "int8->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: int8 = 42;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true, // not allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestInt16TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// int16 -> boolean
		{
			name: "int16->boolean",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 1;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "int16->boolean (other than 0/1)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 1000;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		// int16 -> int8
		{
			name: "int16->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		{
			name: "int16->int8 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 420;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int16 -> int16
		{
			name: "int16->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		// int16 -> int32
		{
			name: "int16->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		// int16 -> int64
		{
			name: "int16->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// int16 -> float16
		{
			name: "int16->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float16, size=4",
		},
		// int16 -> float32
		{
			name: "int16->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float32, size=4",
		},
		// int16 -> float64
		{
			name: "int16->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float64, size=8",
		},
		// int16 -> string
		{
			name: "int16->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: int16 = 42;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true, // not allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestInt32TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// int32 -> boolean
		{
			name: "int32->boolean",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 1;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "int32->boolean (other than 0/1)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 100000;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		// int32 -> int8
		{
			name: "int32->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		{
			name: "int32->int8 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 420;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int32 -> int16
		{
			name: "int32->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		{
			name: "int32->int16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 32769;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int32 -> int32
		{
			name: "int32->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		// int32 -> int64
		{
			name: "int32->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// int32 -> float16
		{
			name: "int32->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float16, size=4",
		},
		{
			name: "int32->float16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 65509;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int32 -> float32
		{
			name: "int32->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float32, size=4",
		},
		// int32 -> float64
		{
			name: "int32->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float64, size=8",
		},
		// int32 -> string
		{
			name: "int32->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: int32 = 42;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true, // not allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestInt64TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// int64 -> boolean
		{
			name: "int64->boolean",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 1;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "int64->boolean (other than 0/1)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 1000000000;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		// int64 -> int8
		{
			name: "int64->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		{
			name: "int64->int8 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 178;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int64 -> int16
		{
			name: "int64->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		{
			name: "int64->int16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 32769;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int64 -> int32
		{
			name: "int64->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		{
			name: "int64->int32 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 2147483649;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int64 -> int64
		{
			name: "int64->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// int64 -> float16
		{
			name: "int64->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float16, size=4",
		},
		{
			name: "int64->float16 (oveflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 65509;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// int64 -> float32
		{
			name: "int64->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float32, size=4",
		},
		// int64 -> float64
		{
			name: "int64->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42.000000, type=float64, size=8",
		},
		// int64 -> string
		{
			name: "int64->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: int64 = 42;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true, // not allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestFloat16TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// float16 -> boolean
		{
			name: "float16->boolean (non-zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 1.5;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "float16->boolean (zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 0.0;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=0, type=boolean, size=1",
		},
		// float16 -> int8
		{
			name: "float16->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 42.0;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		// float16 -> int16
		{
			name: "float16->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 42.0;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		{
			name: "float16->int16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 65504;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float16 -> int32
		{
			name: "float16->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 42.0;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		// float16 -> int64
		{
			name: "float16->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 42.0;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// float16 -> float16
		{
			name: "float16->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 3.14;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140625, type=float16, size=4",
		},
		// float16 -> float32
		{
			name: "float16->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 3.14;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140625, type=float32, size=4",
		},
		// float16 -> float64
		{
			name: "float16->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float16 = 3.14;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140625, type=float64, size=8",
		},
		// float16 -> string
		{
			name: "float16->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: float16 = 3.14;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true, // not allowed
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestFloat32TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// float32 -> boolean
		{
			name: "float32->boolean (non-zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 1.5;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "float32->boolean (zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 0.0;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=0, type=boolean, size=1",
		},
		// float32 -> int8
		{
			name: "float32->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 42.0;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		// float32 -> int16
		{
			name: "float32->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 42.0;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		{
			name: "float32->int16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 1e6;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float32 -> int32
		{
			name: "float32->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 42.0;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		{
			name: "float32->int32 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 1e20;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float32 -> int64
		{
			name: "float32->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 42.0;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		// float32 -> float16
		{
			name: "float32->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 3.14;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			// Note: precision loss
			wantOut: "value=3.140625, type=float16, size=4",
		},
		// float32 -> float32
		{
			name: "float32->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 3.14;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140000, type=float32, size=4",
		},
		// float32 -> float64
		{
			name: "float32->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float32 = 3.14;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140000, type=float64, size=8",
		},
		// float32 -> string (not allowed)
		{
			name: "float32->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: float32 = 3.14;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}

func TestFloat64TypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// float64 -> boolean
		{
			name: "float64->boolean (non-zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 1.5;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=1, type=boolean, size=1",
		},
		{
			name: "float64->boolean (zero)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 0.0;
				say b: boolean = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=0, type=boolean, size=1",
		},
		// float64 -> int8
		{
			name: "float64->int8",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 42.0;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int8, size=1",
		},
		{
			name: "float64->int8 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 1e3;
				say b: int8 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float64 -> int16
		{
			name: "float64->int16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 42.0;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int16, size=2",
		},
		{
			name: "float64->int16 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 1e6;
				say b: int16 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float64 -> int32
		{
			name: "float64->int32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 42.0;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int32, size=4",
		},
		{
			name: "float64->int32 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 1e20;
				say b: int32 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float64 -> int64
		{
			name: "float64->int64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 42.0;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=42, type=int64, size=8",
		},
		{
			name: "float64->int64 (overflow)",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 1e40;
				say b: int64 = a;
				io.printf("value=%d, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantErr: true,
		},
		// float64 -> float16
		{
			name: "float64->float16",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 3.14;
				say b: float16 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140625, type=float16, size=4",
		},
		// float64 -> float32
		{
			name: "float64->float32",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 3.14;
				say b: float32 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140000, type=float32, size=4",
		},
		// float64 -> float64
		{
			name: "float64->float64",
			src: `
			import io from builtin;
			import types from builtin;
			fn main(): int32 {
				say a: float64 = 3.14;
				say b: float64 = a;
				io.printf("value=%f, type=%s, size=%d", b, types.type(b), types.size(b));
				return 0;
			}`,
			wantOut: "value=3.140000, type=float64, size=8",
		},
		// float64 -> string (not allowed)
		{
			name: "float64->string",
			src: `
			import io from builtin;
			fn main(): int32 {
				say a: float64 = 3.14;
				say s: string = a;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := testutils.CompileAndRunSafe(tt.src, t.TempDir())
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("CompileAndRun() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("CompileAndRun() succeeded unexpectedly")
			}
			assert.Equal(t, tt.wantOut, got)
		})
	}
}
