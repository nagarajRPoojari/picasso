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
                say a: int = 100;
                fn Test() {}
                fn tester(do: int, m: Math): Math {
                    if (do == 1){
                        return new Math();
                    }else {
                        return null;
                    }
                    return null;
                }
                fn nuller(): Math {
                    return null;
                }
            }

            class Math {
                say h: int = 90;
                fn Math() {}
            }
			fn main(): int32 {
				say x: Test = new Test();
                say nn: Math;
                say n: Math = x.tester(1, x.nuller());
                io.printf("x= %d", n.h);
				return 0;
			}