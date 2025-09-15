import io;
say PI: double = 3.14;

class DirectoryReader {
  say y: double = 112;
  // say x: Math = new Math();

  // fn sum(m: Math): DirectoryReader {
    
  //   return null;
  // }

  fn math() {
     io.printf("this is inside math");
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

    io.printf("hello world %f\n", 18);

    
    return 0;
}