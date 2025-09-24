package test

import (
	"testing"

	"github.com/nagarajRPoojari/x-lang/utils/testutils"
	"github.com/stretchr/testify/assert"
)

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantOut string
		wantErr bool
	}{
		{
			name: "return int",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn doubleIt(x: int): int {
						return 2 * x;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%d", t.doubleIt(100));
                    return 0;
                }
            `,
			wantOut: "200",
		},
		{
			name: "return float",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn doubleIt(x: float): float {
						return 2 * x;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%f", t.doubleIt(100));
                    return 0;
                }
            `,
			wantOut: "200.000000",
		},
		// returning nothing
		{
			name: "return void by no return statement",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn nothing() {
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
					t.nothing();
                    return 0;
                }
            `,
			wantOut: "",
		},
		{
			name: "return void by empty return statement",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn nothing() {
						return;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
					t.nothing();
                    return 0;
                }
            `,
			wantOut: "",
		},
		{
			name: "void return should not be assignable",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn nothing() {
						return;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
					say n: int = t.nothing();
                    return 0;
                }
            `,
			wantErr: true,
		},
		// return class instance
		{
			name: "return null",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn nothing(): Math {
						return new Math();
                    }
                }
				class Math {
					say x: int = 190;
					fn Math() {}
				}
                fn main(): int32 {
                    say t: Test = new Test();
					say x: Math = t.nothing();
					io.printf("x=%d", x.x);
                    return 0;
                }
            `,
			wantOut: "x=190",
		},
		// null returns
		{
			name: "return null",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn nothing(): Math {
						return null;
                    }
                }
				class Math {
					say x: int = 190;
					fn Math() {}
				}
                fn main(): int32 {
                    say t: Test = new Test();
					say x: Math = t.nothing();
					io.printf("x=%d", x.x);
                    return 0;
                }
            `,
			wantErr: true,
		},
		{
			name: "return by typecasting",
			src: `
                import io;
                class Test {
                    fn Test() {}
                    fn doubleIt(x: int): float {
						return 2 * x;
                    }
                }
                fn main(): int32 {
                    say t: Test = new Test();
                    io.printf("%f", t.doubleIt(100));
                    return 0;
                }
            `,
			wantOut: "200.000000",
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
