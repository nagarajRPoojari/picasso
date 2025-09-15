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

    if n > 10 {
      say y: int = 800;
      if 100 > 20 {
          io.printf("value = %d ", y);
      } else {
          
      }
    } else {
      say ni: int = 200;
    }
    
    return 0;
}