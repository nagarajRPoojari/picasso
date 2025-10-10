// import io from builtin;
// import array from builtin;
// import types from builtin;


// class Any {
//     fn Any() {}
//     fn Print() {
//         io.printf("unimplemented");
//     }
// }

// class Integer: Any {
//     say x: int;
//     fn Integer(x: int) {
//         this.x = x;
//     }

//     fn Print() {
//         io.printf("x=%d. ", this.x);
//     }
// }

// class Test {
//     say x: [][]Integer = array.create(Integer, 2,2);
//     fn Test() {}
//     fn Update(arr: [][]Integer): [][]Integer {
//         arr[0,1].Print();
//         say n: Integer = arr[0,1];
//         n.x = 90;
//         return arr;
//     }   
// }
// fn main(): int32 {
//     say arr: [][]Integer;
//     arr = array.create(Integer, 5, 4);
//     arr[0, 1] = new Integer(100);

//     say t: Test = new Test();
//     say x: [][]Integer = t.Update(arr);

//     t.x = arr;

//     t.Update(arr);
//     return 0;
// }

                // import io from builtin;
                // import array from builtin;
                // class Test {
                //   say x: [][]Integer;
                //   fn Test() {}
                // }
                // class Integer {
                //     say x: int;
                //     fn Integer(v: int) {
                //         this.x = v;
                //     }
                // }
                // fn main(): int32 {
                //     say c: Test = new Test();
                //     c.x = array.create(Integer, 2, 2);
                //     c.x[1,1] = new Integer(20);
                //     io.printf("x=%d", c.x[1,1].x);
                //     return 0;
                // }


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