package test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeclareVar(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// declarations with literals
		{
			name: "int declaration",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 10;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "10",
		},
		{
			name: "float declaration",
			src: `
                import io;
                fn main(): int32 {
                    say f: float64 = 3.14159;
                    io.printf("%f", f);
                    return 0;
                }
            `,
			wantOut: "3.141590",
		},
		{
			name: "string declaration",
			src: `
                import io;
                fn main(): int32 {
                    say s: string = "hello";
                    io.printf("%s", s);
                    return 0;
                }
            `,
			wantOut: "hello",
		},
		{
			name: "boolean declaration",
			src: `
                import io;
                fn main(): int32 {
                    say s: boolean = 1;
                    io.printf("%d", s);
                    return 0;
                }
            `,
			wantOut: "1",
		},
		{
			name: "declaration as class fields",
			src: `
                import io;
                class Test {
                    say x: int = 42;
					fn Test() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d", t.x);
                    return 0;
                }
            `,
			wantOut: "42",
		},
		// declaration with expression
		{
			name: "init with expression",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 5;
                    say b: int = a + 10;
                    io.printf("%d %d", a, b);
                    return 0;
                }
            `,
			wantOut: "5 15",
		},
		{
			name: "literals with unary expression",
			src: `
                import io;
                fn main(): int32 {
                    say z: int = 0;
                    say n: int = -100;
                    io.printf("%d %d", z, n);
                    return 0;
                }
            `,
			wantOut: "0 -100",
		},
		{
			name: "init class fields with expression",
			src: `
                import io;
                class Test {
                    say x: int = 42;
					say y: int = 2 * x;
					fn Test() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d", t.y);
                    return 0;
                }
            `,
			wantOut: "84",
		},
		// uninitialized variables
		{
			name: "uninitialized int, float, boolean",
			src: `
                import io;
                fn main(): int32 {
                    say a: int;
					say b: float;
					say c: boolean;
                    io.printf("%d %f %d", a, b, c);
                    return 0;
                }
            `,
			wantOut: "0 0.000000 0",
		},
		{
			name: "uninitialized string",
			src: `
                import io;
                fn main(): int32 {
                    say a: string;
                    io.printf("%s", a);
                    return 0;
                }
            `,
			wantOut: "(null)",
		},
		{
			name: "uninitialized class fields",
			src: `
                import io;
                class Test {
                    say x: int;
					fn Test() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d", t.x);
                    return 0;
                }
            `,
			wantOut: "0",
		},
		// redeclare should fail
		{
			name: "redeclare variable",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 10;
                    say a: int = 20;
                    return 0;
                }
            `,
			wantErr: true,
		},
		{
			name: "redeclaring class fields",
			src: `
                import io;
                class Test {
                    say x: int;
					say x: int;
					fn Test() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d", t.x);
                    return 0;
                }
            `,
			wantErr: true,
		},
		// visibility resolution
		{
			name: "different funcs with same var names",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn fa(): int {
                      say a: int = 100;
                      return a;
                    }
              
                    fn fb(): int {
                      say a: int = 200;
                      return a;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d %d", t.fa(), t.fb());
                    return 0;
                }
            `,
			wantOut: "100 200",
		},
		{
			name: "name resolution in nested blocks",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 100;
                    io.printf("a=%d", a);
                    if a > 10 {
                      say a: int = 200;
                      io.printf("a=%d", a);
                    }else {
                      io.printf("%d is less than 10", a);
                    }
                    return 0;
                }
            `,
			wantOut: "a=100a=200",
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

func TestAssignVar(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "int reassignment",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 10;
					a = 180;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "180",
		},
		{
			name: "float reassignment with int literal",
			src: `
                import io;
                fn main(): int32 {
                    say f: float64 = 3.14159;
					f = 90;
                    io.printf("%f", f);
                    return 0;
                }
            `,
			wantOut: "90.000000",
		},
		{
			name: "string reassignment",
			src: `
                import io;
                fn main(): int32 {
                    say s: string = "hello";
					s = "never";
                    io.printf("%s", s);
                    return 0;
                }
            `,
			wantOut: "never",
		},
		// assignment with expression
		{
			name: "init with expression",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 90;
                    a = a + 10;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "100",
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
