// import io;
// say PI: double = 3.14;

// class DirectoryReader {
//   say y: double = 112;
//   say x: Math = new Math();
//   say s: string = "nagaraj";

//   fn DirectoryReader() {
//     // this.y = 100;
//     // this.x.pi = 98;
//     // io.printf("this.x.y = %f          ", this.y);
//   }

//   fn math(intern: Math, inn: string): Math {
//     say a: Math = new Math();
//     io.printf("this.x.pi = %f  %f %f  %s, ", this.x.pi, a.pi, intern.pi, inn);
//     // this.add();
//     this.x.greet();
//     return a;
//   }

//   fn add(): double {
//     // say a: Math = new Math();
//      return 90.0;
//   }

//   fn str(): string {
//     return this.s;
//   }
// }

// class Math {
//   say pi: double = 123;
//   fn Math() {
//     // io.printf("printing pi:  %f -- ", this.pi);
//   } 

//   fn greet() {
//     // io.printf("I am inside greet method.   ");
//   }
// }

// fn main(): int32 {
//     say a: DirectoryReader = new DirectoryReader();
//     say n: int;
//     n = 278;
//     say m: Math = a.math(new Math(), "hosad");
  
//     m.greet();
//     say z: string = a.str();
//     io.printf("last:   %f   ", m.pi);
//     if a.add()  > 10 {
//       say y: int = 800;
//       if 100 > 20 {
//           io.printf("value = %d ", y);
//       } else {
          
//       }
//     } else {
//       say ni: int = 200;
//     }
    
//     return 0;
// }
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