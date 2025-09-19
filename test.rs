import io;
say PI: double = 3.14;

class DirectoryReader {
  say y: double = 112;
  say x: Math = new Math();

  fn DirectoryReader() {
    // this.y = 100;
    // this.x.pi = 98;
    // io.printf("this.x.y = %f          ", this.y);
  }

  fn math(intern: Math): Math {
    say a: Math = new Math();
    io.printf("this.x.pi = %f  %f %f  , ", this.x.pi, a.pi, intern.pi);
    // this.add();
    // this.x.greet();
    return a;
  }

  fn add(): double {
    //  say a: Math = new Math();
     return this.y;
  }

  fn str(): string {
    return "hello";
  }
}

class Math {
  say pi: double = 123;
  fn Math() {
    // io.printf("printing pi:  %f -- ", this.pi);
  } 

  fn greet() {
    // io.printf("I am inside greet method.   ");
  }
}

fn main(): int32 {
    say a: DirectoryReader = new DirectoryReader();
    say n: int;
    n = 278;
    say m: Math = a.math(new Math());
  
    m.greet();
    say z: string = a.str();
    io.printf("last:   %f   ", m.pi);
    // if a.add()  > 10 {
    //   say y: int = 800;
    //   if 100 > 20 {
    //       io.printf("value = %d ", y);
    //   } else {
          
    //   }
    // } else {
    //   say ni: int = 200;
    // }
    
    return 0;
}