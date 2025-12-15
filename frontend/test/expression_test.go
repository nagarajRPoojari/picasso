package test

import (
	"testing"

	"github.com/nagarajRPoojari/niyama/frontend/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDeclareVarExpression(t *testing.T) {
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
                fn main(): int32 {
                    say s: boolean = 1;
                    io.printf("%d", s);
                    return 0;
                }
            `,
			wantOut: "1",
		},
		{
			name: "declaration as class instance",
			src: `
                import io from builtin;
                class Test {
                    say x: Math = new Math();
					fn Test() {}
                }
                class Math {
                    fn Math() {}
                } 
                fn main(): int32 {
                    say t: Test = new Test();
                    return 0;
                }
            `,
			wantOut: "",
		},
		{
			name: "declaration as class fields",
			src: `
                import io from builtin;
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
		{
			name: "declaration of array as class fields",
			src: `
                import io from builtin;
                import array from builtin;
                class Test {
                    say x: [][]Math = array.create(Math, 3, 4);
					fn Test() {
                        this.x[0,0] = new Math();
                    }
                }
                class Math {
                    say pi: double = 3.14;
                    fn Math() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%f", t.x[0,0].pi);
                    return 0;
                }
            `,
			wantOut: "3.140000",
		},
		{
			name: "2D array of class types declaration",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;
                class Test {
                    say x: int = 42;
					fn Test() {}
                }
                fn main(): int32 {
                    say arr: [][]Test = array.create(Test, 4, 5);
                    arr[0,0] = new Test();
                    io.printf("type=%s,eleType=%s,eleValue at [0,0]=%d", types.type(arr), types.type(arr[0,0]), arr[0,0].x);
                    return 0;
                }
            `,
			wantOut: "type=array,eleType=Test,eleValue at [0,0]=42",
		},
		// declaration with expression
		{
			name: "init with expression",
			src: `
                import io from builtin;
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
                import io from builtin;
                fn main(): int32 {
                    say z: int = 0;
                    say n: int = -100;
                    io.printf("%d %d", z, n);
                    return 0;
                }
            `,
			wantOut: "0 -100",
		},
		// uninitialized variables
		{
			name: "uninitialized int, float, boolean",
			src: `
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
		{
			name: "uninitialized class instance type",
			src: `
                import io from builtin;
                class Test {
                    say x: Math;
					fn Test() {}
                }
                class Math {
                    fn Math() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    return 0;
                }
            `,
			wantOut: "",
		},
		{
			name: "uninitialized array type",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;
                fn main(): int32 {
                    say arr: []int;
                    io.printf("type = %s ", types.type(arr));
                    return 0;
                }
            `,
			wantOut: "type = array ",
		},
		{
			name: "allow recursive type reference",
			src: `
                import io from builtin;
                class Test {
                    say x: Test;
					fn Test() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    return 0;
                }
            `,
			wantOut: "",
		},
		// redeclare should fail
		{
			name: "redeclare variable",
			src: `
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
		// type cast
		{
			name: "type cast float64 to int8",
			src: `
                import io from builtin;
                fn main(): int32 {
                    say a: int8 = 1899;
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

func TestAssignVarExpression(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "int reassignment",
			src: `
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
                fn main(): int32 {
                    say s: string = "hello";
					s = "never";
                    io.printf("%s", s);
                    return 0;
                }
            `,
			wantOut: "never",
		},
		{
			name: "class type reassignment",
			src: `
                import io from builtin;
                class Test {
                    say x: Math;
                    fn Test() {
                        this.x = new Math();
                    }
                }
                class Math {
                    say x: int = 100;
                    fn Math() {}
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("x=%d", t.x.x);
                    return 0;
                }
            `,
			wantOut: "x=100",
		},
		{
			name: "array type reassignment",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;
                fn main(): int32 {
                    say arr: []int = array.create(int, 10);
                    arr = array.create(int, 2);
                    io.printf("type = %s ", types.type(arr));
                    return 0;
                }
            `,
			wantOut: "type = array ",
		},
		// assignment with expression
		{
			name: "init with expression",
			src: `
                import io from builtin;
                fn main(): int32 {
                    say a: int = 90;
                    a = a + 10;
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "100",
		},
		{
			name: "init with member expression",
			src: `
                import io from builtin;
                class Test {
                    say x: int = 100;
                    fn Test() {}
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    say b: int = a.x;
                    io.printf("%d", b);
                    return 0;
                }
            `,
			wantOut: "100",
		},
		{
			name: "init with array index expression",
			src: `
                import io from builtin;
                import array from builtin;
                fn main(): int32 {
                    say arr: []int = array.create(int, 4);
                    say b: int = arr[0];
                    io.printf("%d", b);
                    return 0;
                }
            `,
			wantOut: "0",
		},
		// initializing the uninitialized
		{
			name: "uninitialized array assignment",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;
                fn main(): int32 {
                    say arr: []int;
                    arr = array.create(int, 2);
                    arr[0] = 10;
                    io.printf("type = %s, %d ", types.type(arr), arr[0]);
                    return 0;
                }
            `,
			wantOut: "type = array, 10 ",
		},
		{
			name: "uninitialized class assignment",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {this.x = x;}
                }
                fn main(): int32 {
                    say t: Test;
                    t = new Test(90);
                    io.printf("type = %s, %d ", types.type(t), t.x);
                    return 0;
                }
            `,
			wantOut: "type = Test, 90 ",
		},
		// visibility
		{
			name: "name resolution with function params",
			src: `
                import io from builtin;
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
                import io from builtin;
                class Test {
                  say x: float;
                  fn Test() {}
                }

                fn main(): int32 {
                    say c: Test = new Test();
                    c.x = 190;
                    io.printf("x=%f", c.x);
                    return 0;
                }
            `,
			wantOut: "x=190.000000",
		},
		{
			name: "assign class fields: array",
			src: `
                import io from builtin;
                import array from builtin;
                class Test {
                  say x: [][]Integer;
                  fn Test() {}
                }
                class Integer {
                    say x: int;
                    fn Integer(v: int) {
                        this.x = v;
                    }
                }
                fn main(): int32 {
                    say c: Test = new Test();
                    c.x = array.create(Integer, 2, 2);
                    c.x[1,1] = new Integer(20);
                    io.printf("x=%d", c.x[1,1].x);
                    return 0;
                }
            `,
			wantOut: "x=20",
		},
		// invalid assignments
		{
			name: "name resolution with function params",
			src: `
                import io from builtin;
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
                import io from builtin;
                fn main(): int32 {
                    say a: int = 100;
					a = "hello";
                    return 0;
                }
            `,
			wantErr: true,
		},
		// type cast
		{
			name: "type cast float64 to int8",
			src: `
                import io from builtin;
                fn main(): int32 {
                    say a: int8 = 1090;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
                import array from builtin;
                fn main(): int32 {
                    say x: []int = array.create(int, 2);
                    x[0] = 100;
                    say a: int = 2 + 3 * 4 * x[0];
                    io.printf("%d", a);
                    return 0;
                }
            `,
			wantOut: "1202",
		},
		{
			name: "Parentheses override precedence",
			src: `
                import io from builtin;
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
                import io from builtin;
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

func TestFunctionCallExpression(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "function call without param, without return",
			src: `
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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
                import io from builtin;
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

		// function call involving heap vars
		{
			name: "class types as function params",
			src: `
                import io from builtin;
                class Test {
                    fn Test() {}
                    fn Calc(m: Math): double {
                        say prev:double =  m.PI;
                        m.Reset();
                        return prev;
                    }
                }
                class Math {
                    say PI: double = 3.14;
                    fn Math() {}
                    fn Reset() {
                        this.PI = 0.0;
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    say m: Math = new Math();
                    io.printf("before: %f", a.Calc(m));
                    io.printf("after: %f", m.PI);
                    return 0;
                }
            `,
			wantOut: "before: 3.140000after: 0.000000",
		},
		{
			name: "array types as function params",
			src: `
                import io from builtin;
                import array from builtin;
                import types from builtin;


                class Any {
                    fn Any() {}
                    fn Print() {
                        io.printf("unimplemented");
                    }
                }

                class Integer: Any {
                    say x: int;
                    fn Integer(x: int) {
                        this.x = x;
                    }

                    fn Print() {
                        io.printf("x=%d. ", this.x);
                    }
                }

                class Test {
                    say x: [][]Integer = array.create(Integer, 2,2);
                    fn Test() {}
                    fn Update(arr: [][]Integer): [][]Integer {
                        arr[0,1].Print();
                        say n: Integer = arr[0,1];
                        n.x = 90;
                        return arr;
                    }   
                }
                fn main(): int32 {
                    say arr: [][]Integer;
                    arr = array.create(Integer, 5, 4);
                    arr[0, 1] = new Integer(100);

                    say t: Test = new Test();
                    say x: [][]Integer = t.Update(arr);

                    t.x = arr;

                    t.Update(arr);
                    return 0;
                }
            `,
			wantOut: "x=100. x=90. ",
		},
		// returns
		{
			name: "return class types",
			src: `
                import io from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {
                        this.x = x;
                    }
                    fn self(): Test {
                        return this;
                    }
                }

                fn main(): int32 {
                    say a: Test = new Test(89);  
                    say b: Test = a.self();
                    io.printf("%d", b.x);
                    return 0;
                }
            `,
			wantOut: "89",
		},
		{
			name: "return array types",
			src: `
                import io from builtin;
                import array from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {
                        this.x = x;
                    }
                    fn self(): [][]Test {
                        say x: [][]Test = array.create(Test, 2,2);
                        x[1,1] = new Test(78);
                        return x;
                    }
                }

                fn main(): int32 {
                    say a: Test = new Test(89);  
                    say b: [][]Test = a.self();
                    io.printf("%d", b[1,1].x);
                    return 0;
                }
            `,
			wantOut: "78",
		},
		// type cast
		{
			name: "type cast function return",
			src: `
                import io from builtin;
                class Test {
                    fn Test() {}
                    fn printer(x: float): float {
                        return x;
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();  
                    say b: int = a.printer(100);
                    io.printf("%d", b);
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

func TestNewExpression(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "basic constructor with field modification",
			src: `
                import io from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {
                        this.x = x;
                    }
                    fn greet() {
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test(200);
                    io.printf("%d", t.x);
                    return 0;
                }
            `,
			wantOut: "200",
		},
		{
			name: "independent new expression",
			src: `
                import io from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {
                        this.x = x;
                    }
                    fn greet() {
                    }
                }
                fn main(): int32 {
                    new Test(200);
                    return 0;
                }
            `,
			wantOut: "",
		},
		{
			name: "constructor as a method",
			src: `
                import io from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int) {
                        this.x = x;
                    }
                    fn increment() {
                        this.Test(this.x + 1);
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test(200);
                    t.increment();

                    io.printf("after: x=%d", t.x);
                    return 0;
                }
            `,
			wantOut: "after: x=201",
		},
		{
			name: "constructor as a method with returns",
			src: `
                import io from builtin;
                class Test {
                    say x: int;
                    fn Test(x: int): Test {
                        this.x = x;
                        return this;
                    }
                    fn increment(): Test {
                        return this.Test(this.x + 1);
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test(200);
                    say x:Test = t.increment();

                    io.printf("after: x=%d", x.x);
                    return 0;
                }
            `,
			wantOut: "after: x=201",
		},
		{
			name: "new expression involving array as param",
			src: `
                import io from builtin;
                import array from builtin;
                class Test {
                    say x: []int;
                    fn Test(x: []int) {
                        this.x = x;
                    }
                }
                fn main(): int32 {
                    say t:Test = new Test(array.create(int, 2));
                    return 0;
                }
            `,
			wantOut: "",
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
