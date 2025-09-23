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
		{
			name: "name resolution with function params",
			src: `
                import io;
                class Test {
                  fn Test() {
                  }

                  fn test(x: int): int {
				    // should be able to hide params as well
                    say x: int = 20;
                    return x;
                  }

                }

                fn main(): int32 {
                    say a: int = 100;
                    say c: Test = new Test();
                    io.printf("a=%d", c.test(1000));
                    return 0;
                }
            `,
			wantOut: "a=20",
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
		// visibility
		{
			name: "name resolution with function params",
			src: `
                import io;
                class Test {
                  fn Test() {
                  }

                  fn test(x: int): int {
				    // should be able to hide params as well
                    x = 20;
                    return x;
                  }

                }

                fn main(): int32 {
                    say a: int = 100;
                    say c: Test = new Test();
                    io.printf("a=%d", c.test(1000));
                    return 0;
                }
            `,
			wantOut: "a=20",
		},
		// assign class fields
		{
			name: "assign class fields",
			src: `
                import io;
                class Test {
                  say x: float;
                  fn Test() {
                  }
                }

                fn main(): int32 {
                    say c: Test = new Test();
                    c.x = 190;
                    io.printf("x=%d", x);
                    return 0;
                }
            `,
			wantOut: "x=190",
		},
		// invalid assignments
		{
			name: "name resolution with function params",
			src: `
                import io;
                class Test {
                  fn Test() {
                  }

                  fn test() {
                  }

                }

                fn main(): int32 {
                    say a: int = 100;
                    say c: Test = new Test();
					a = c.test();
                    return 0;
                }
            `,
			wantErr: true,
		},
		{
			name: "name resolution with function params",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 100;
					a = "hello";
                    return 0;
                }
            `,
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

func TestBinaryExpression(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		// arithmetic
		{
			name: "Addition",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 5 + 7;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "12",
		},
		{
			name: "Subtraction with negative result",
			src: `
                import io;
                fn main(): int32 {
                    say b: int = 10 - 20;
                    io.printf("%d", b);
                    return 0;
                }
            `,
			wantOut: "-10",
		},
		{
			name: "Multiplication",
			src: `
                import io;
                fn main(): int32 {
                    say c: int = 6 * 7;
                    io.printf("%d", c);
                    return 0;
                }
            `,
			wantOut: "42",
		},
		{
			name: "Division",
			src: `
                import io;
                fn main(): int32 {
                    say d: int = 20 / 4;
                    io.printf("%d", d);
                    return 0;
                }
            `,
			wantOut: "5",
		},
		// comparision
		{
			name: "Greater than true",
			src: `
                import io;
                fn main(): int32 {
                    say x: boolean = 10 > 5;
                    io.printf("%d", x);
                    return 0;
                }
            `,
			wantOut: "1",
		},
		{
			name: "Less than false",
			src: `
                import io;
                fn main(): int32 {
                    say y: boolean = 10 < 5;
                    io.printf("%d", y);
                    return 0;
                }
            `,
			wantOut: "0",
		},
		{
			name: "Equality",
			src: `
                import io;
                fn main(): int32 {
                    say z: boolean = 7 == 7;
                    io.printf("%d", z);
                    return 0;
                }
            `,
			wantOut: "1",
		},
		{
			name: "Inequality",
			src: `
                import io;
                fn main(): int32 {
                    say z: boolean = 7 != 8;
                    io.printf("%d", z);
                    return 0;
                }
            `,
			wantOut: "1",
		},

		// logical
		{
			name: "Logical AND",
			src: `
                import io;
                fn main(): int32 {
                    say res: boolean = 1 && 0;
                    io.printf("%d", res);
                    return 0;
                }
            `,
			wantOut: "0",
		},
		{
			name: "Logical OR",
			src: `
                import io;
                fn main(): int32 {
                    say res: boolean = 0 || 1;
                    io.printf("%d", res);
                    return 0;
                }
            `,
			wantOut: "1",
		},
		{
			name: "Logical NOT",
			src: `
                import io;
                fn main(): int32 {
                    say res: boolean = !0;
                    io.printf("%d", res);
                    return 0;
                }
            `,
			wantOut: "1",
		},

		// mixed precedence
		{
			name: "Precedence multiplication before addition",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = 2 + 3 * 4;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "14",
		},
		{
			name: "Parentheses override precedence",
			src: `
                import io;
                fn main(): int32 {
                    say a: int = (2 + 3) * 4;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "20",
		},
		{
			name: "involving function calls",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn pi(): float {
                        return 3.14;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    say a: int = (2 + 3) * 4 * t.pi();
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "62",
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

func TestCallFunc(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "function call without param, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn greet() {
                        io.printf("hi");
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.greet();
                    return 0;
                }
            `,
			wantOut: "hi",
		},
		{
			name: "function call with params, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000",
		},
		{
			name: "function call with params, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000",
		},
		{
			name: "function call with params, with return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double): string {
                        io.printf("x=%d, y=%f", x, y);
                        return "hello";
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    io.printf("%s", a.printer(10,20));
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000hello",
		},
		{
			name: "unknown function call",
			src: `
                import io;
                class Test {
                    fn Test() {}
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.greet();
                    return 0;
                }
            `,
			wantErr: true,
		},
		{
			name: "function call with ignoring params",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    a.printer(10);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000x=10, y=0.000000",
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

func TestCallConstructor(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "function call without param, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn greet() {
                        io.printf("hi");
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.greet();
                    return 0;
                }
            `,
			wantOut: "hi",
		},
		{
			name: "function call with params, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000",
		},
		{
			name: "function call with params, without return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000",
		},
		{
			name: "function call with params, with return",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double): string {
                        io.printf("x=%d, y=%f", x, y);
                        return "hello";
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    io.printf("%s", a.printer(10,20));
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000hello",
		},
		{
			name: "unknown function call",
			src: `
                import io;
                class Test {
                    fn Test() {}
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.greet();
                    return 0;
                }
            `,
			wantErr: true,
		},
		{
			name: "function call with ignoring params",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: double) {
                        io.printf("x=%d, y=%f", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,20);
                    a.printer(10);
                    return 0;
                }
            `,
			wantOut: "x=10, y=20.000000x=10, y=0.000000",
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
