package test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestAllTypeCasting(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// boolean casts
		{
			name: "boolean->boolean",
			src: `
			import io;
			fn main(): int32 {
				say b: boolean = 1;
				say c: boolean = b;
				io.printf("%d", c);
				return 0;
			}`,
			wantOut: "1",
		},
		// {
		// 	name: "boolean->int",
		// 	src: `
		// 	import io;
		// 	fn main(): int32 {
		// 		say b: boolean = 1;
		// 		say x: int = b;
		// 		io.printf("%d", x);
		// 		return 0;
		// 	}`,
		// 	wantOut: "1",
		// },
		// {
		// 	name: "boolean->int8",
		// 	src: `
		// 	import io;
		// 	fn main(): int32 {
		// 		say b: boolean = false;
		// 		say x: int8 = b;
		// 		io.printf("%d", x);
		// 		return 0;
		// 	}`,
		// 	wantOut: "0",
		// },
		// {
		// 	name: "boolean->float64",
		// 	src: `
		// 	import io;
		// 	fn main(): int32 {
		// 		say b: boolean = true;
		// 		say f: float64 = b;
		// 		io.printf("%f", f);
		// 		return 0;
		// 	}`,
		// 	wantOut: "1.000000",
		// },
		{
			name: "boolean->string",
			src: `
			import io;
			fn main(): int32 {
				say b: boolean = true;
				say s: string = b;
				io.printf("%s", s);
				return 0;
			}`,
			wantErr: true}, // not allowed

		// int casts
		{
			name: "int->int8 (truncate)",
			src: `
			import io;
			fn main(): int32 {
				say x: int = 257;
				say y: int8 = x;
				io.printf("%d", y);
				return 0;
			}`,
			wantOut: "1",
		},
		{
			name: "int->int16",
			src: `
			import io;
			fn main(): int32 {
				say x: int = 32001;
				say y: int16 = x;
				io.printf("%d", y);
				return 0;
			}`,
			wantOut: "-335",
		},
		{
			name: "int->int64",
			src: `
			import io;
			fn main(): int32 {
				say x: int = 123456789;
				say y: int64 = x;
				io.printf("%d", y);
				return 0;
			}`,
			wantOut: "123456789",
		},
		// float casts
		{
			name: "int32->float32",
			src: `
			import io;
			fn main(): int32 {
				say x: int32 = 42;
				say f: float32 = x;
				io.printf("%f", f);
				return 0;
			}`,
			wantOut: "42.000000",
		},
		{
			name: "int64->float64",
			src: `
			import io;
			fn main(): int32 {
				say x: int64 = 123456789;
				say f: float64 = x;
				io.printf("%f", f);
				return 0;
			}`,
			wantOut: "123456789.000000",
		},
		// float to int
		{
			name: "float64->int32 (truncate)",
			src: `
			import io;
			fn main(): int32 {
				say f: float64 = 3.99;
				say i: int32 = f;
				io.printf("%d", i);
				return 0;
			}`,
			wantOut: "3",
		},
		{
			name: "double->int8 (truncate & wrap)",
			src: `
			import io;
			fn main(): int32 {
				say f: double = 257.99;
				say i: int8 = f;
				io.printf("%d", i);
				return 0;
			}`,
			wantOut: "1",
		},
		// float to float
		{
			name: "float16->float32",
			src: `  
			import io;
			fn main(): int32 {
				say f16: float16 = 1.5;
				say f32: float32 = f16;
				io.printf("%f", f32);
				return 0;
			}`,
			wantOut: "1.500000",
		},
		{
			name: "float32->float64",
			src: ` 
			import io;
			fn main(): int32 {
				say f32: float32 = 3.14159;
				say f64: float64 = f32;
				io.printf("%f", f64);
				return 0;
			}`,
			wantOut: "3.141590",
		},
		{
			name: "float64->double",
			src: ` 
			import io;
			fn main(): int32 {
				say f64: float64 = 2.71828;
				say d: double = f64;
				io.printf("%f", d);
				return 0;
			}`,
			wantOut: "2.718280",
		},

		{
			name: "double->float32 (lossy)",
			src: `
			import io;
			fn main(): int32 {
				say d: double = 1.23456789;
				say f: float32 = d;
				io.printf("%f", f);
				return 0;
			}`,
			wantOut: "1.234568"},

		// int/float to string
		{
			name: "int->string (invalid)",
			src: `
			import io;
			fn main(): int32 {
				say x: int = 42;
				say s: string = x;
				return 0;
			}`,
			wantErr: true,
		},
		{
			name: "float64->string (invalid)",
			src: `
			import io;
			fn main(): int32 {
				say f: float64 = 3.14;
				say s: string = f;
				return 0;
			}`,
			wantErr: true,
		},

		// string to other types
		{
			name: "string->int (invalid)",
			src: `
			import io;
			fn main(): int32 {
				say s: string = "123";
				say x: int = s;
				return 0;
			}`,
			wantErr: true,
		},

		{
			name: "string->boolean (invalid)",
			src: `
			import io;
			fn main(): int32 {
				say s: string = "true";
				say b: boolean = s;
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
