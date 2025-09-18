import io;
say PI: double = 3.14;

class DirectoryReader {
  say y: double = 112;
  say x: Math = new Math();

  fn DirectoryReader() {
    this.y = 100;
    this.x.pi = 98;
}

  fn math() {
     io.printf("this.y = %f  ", this.x.pi);
     this.add();

  }

  fn add() {
     io.printf("this is inside math");
     
  }
}

class Math {
  say pi: double = 123;
  fn Math() {

  }
}

fn main(): int32 {
   say arr: [10][2]string = [["1,", "2", "3"],["1,", "2", "3"]];
  //  say xxx: string = arr[0][1][3];

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