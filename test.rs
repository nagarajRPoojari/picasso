import io;
say PI: double = 3.14;

class DirectoryReader {
  say y: double = 112;
  // say x: Math = new Math();

  fn math() {
    //  io.printf("this is inside math");
  }
}

class Math {
  say pi: double = 123;
}

fn main(): int32 {
    say a: DirectoryReader = new DirectoryReader();
    say n: int;
    n = 278;
    a.math();
    say z: string = "hello world";

    io.printf("hello world %s \n", z);

    // if n > 10 {
    //   io.printf("hello world %s \n", z);
    // } else {
    //   say n: int = 200;
    // }
    
    return 0;
}