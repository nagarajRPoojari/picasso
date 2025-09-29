                import io from builtin;
                class Test {
                    fn Test() {}
                    fn printer(x: int, y: string) {
                        io.printf("x=%d, y=%s", x, y);
                    }
                }
                fn main(): int32 {
                    say a: Test = new Test();
                    a.printer(10,"hi");
                    a.printer(10);
                    return 0;
                }