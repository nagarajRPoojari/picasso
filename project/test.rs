                import io from builtin;
                import array from builtin;
                import types from builtin;
                class Test {
                    say x: int = 90;
                    fn Test(x: int){this.x = x;}
                }
                fn main(): int32 {
                    say arr: []int;
                    arr = array.create(int, 2);

                    arr[0,0] = 90;

                    io.printf("type = %d ", arr[0,0]);
                    // say c: Test;
                    // c = new Test(80);

                    // io.printf("printing: %d ", c.x);

                    return 0;
                }